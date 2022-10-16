package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"wbintern/l0/service/cache"
	"wbintern/l0/service/config"
	"wbintern/l0/service/database"
	"wbintern/l0/service/model"
	"wbintern/l0/service/server"

	"github.com/nats-io/stan.go"
)

const usage = `Usage: service [FILE]`

var (
	ch cache.Cache
	db *database.Database
)

func main() {
	// Get config filename as an argument
	cfgpath := "config.toml"

	if len(os.Args) > 1 {
		cfgpath = os.Args[1]
	}

	if _, err := os.Stat(cfgpath); errors.Is(err, os.ErrNotExist) {
		log.Fatal(err, "\n", usage)
	}

	// Parse config
	var cfg, err = config.ParseConfig(cfgpath)
	checkError(err)

	// Establish DB connection
	connStr :=
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.Database.User, cfg.Database.Pass, cfg.Database.Host, cfg.Database.Port, cfg.Database.Db)
	db, err = database.Connect(connStr)
	checkError(err)
	defer db.Close(context.Background())

	// Initialize cache store
	ch = cache.New()
	err = ch.Restore(db)
	checkError(err)

	// Subscribe to nats streaming channel
	natsUrl := fmt.Sprintf("%s:%d", cfg.Stan.Host, cfg.Stan.Port)
	conn, err := stan.Connect(cfg.Stan.Clusterid, cfg.Stan.Userid, stan.NatsURL(natsUrl))
	checkError(err)
	defer conn.Close()

	sub, err := conn.Subscribe(cfg.Stan.Channel, messageHandler)
	checkError(err)
	defer sub.Close()

	// Start HTTP server
	address := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	fmt.Println(address)
	srv := server.New(address, ch)
	srv.ListenAndServe()

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Handle json-formated order from channel
// and store it in a cache and database
func messageHandler(msg *stan.Msg) {
	// Add to cache
	data, err := model.NewFromByte(msg.Data)
	if err != nil {
		log.Println("Can't cache message: ", err)
		return
	}
	ch.Add(*data)

	err = db.InsertOrder(*data)
	if err != nil {
		log.Println("Can't insert data into database: ", err)
		return
	}
}
