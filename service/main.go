package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"wbintern/l0/service/config"
	"wbintern/l0/service/database"
	"wbintern/l0/service/model"

	"github.com/nats-io/stan.go"
)

const usage = `Usage: service [FILE]`

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
	// Subscribe to nats streaming channel
	conn, err := stan.Connect(cfg.Stan.Clusterid, cfg.Stan.Userid, stan.NatsURL(cfg.Stan.Host))
	checkError(err)
	defer conn.Close()

	sub, err := conn.Subscribe(cfg.Stan.Channel, messageHandler, stan.StartWithLastReceived())
	checkError(err)
	defer sub.Close()

	// Establish DB connection
	connStr := fmt.Sprintf("postgres://%v:%v@%v/%v", cfg.Database.User, cfg.Database.Pass, cfg.Database.Host, cfg.Database.Db)
	db, err := database.Connect(connStr)
	checkError(err)
	defer db.Connection.Close(context.Background())

	// Initialize cache store

	// Start HTTP server
	http.ListenAndServe(":8080", nil)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// temp
func messageHandler(msg *stan.Msg) {
	data, err := model.NewFromByte(msg.Data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(*data)
}
