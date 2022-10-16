package app

import (
	"context"
	"fmt"
	"log"
	"wbintern/l0/service/cache"
	"wbintern/l0/service/config"
	"wbintern/l0/service/database"
	"wbintern/l0/service/model"
	"wbintern/l0/service/server"

	"github.com/nats-io/stan.go"
)

type App struct {
	cfg    config.Config
	ch     cache.Cache
	db     *database.Database
	stan   *stan.Conn
	server *server.Server
}

func New(cfg config.Config) (*App, error) {
	log.Println("App Initialization")
	log.Println("Database Connection")
	connStr :=
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.Database.User, cfg.Database.Pass, cfg.Database.Host, cfg.Database.Port, cfg.Database.Db)
	log.Printf("Connecting to %s:%d...", cfg.Database.Host, cfg.Database.Port)
	db, err := database.Connect(connStr)
	if err != nil {
		return &App{}, err
	}
	log.Println("Database Connection:\tSUCCESSFUL")

	log.Println("Cache Initialization")
	ch := cache.New()
	log.Println("Restoring cache from database...")
	err = ch.Restore(db)
	if err != nil {
		return &App{}, err
	}
	log.Println("Cache Initialization:\tSUCCESSFUL")

	log.Println("Stan Connection")
	connStr = fmt.Sprintf("%s:%d", cfg.Stan.Host, cfg.Stan.Port)
	log.Printf("Connecting to %s...\n", connStr)
	stan, err := stan.Connect(cfg.Stan.Clusterid, cfg.Stan.Userid, stan.NatsURL(connStr))
	if err != nil {
		return &App{}, err
	}
	log.Println("Stan Connection:\t\tSUCCESSFUL")

	log.Println("Creating Server")
	connStr = fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := server.New(connStr, ch)
	log.Println("Server Created")

	log.Println("App Initialization:\t\tSUCCESSFUL")
	return &App{
		cfg:    cfg,
		ch:     ch,
		db:     db,
		stan:   &stan,
		server: srv,
	}, nil
}

func (a *App) Close() {
	a.db.Close(context.Background())
	(*a.stan).Close()
	a.server.Close()
	a.ch = nil
}

func (a *App) Run() error {
	log.Println("Stan channel subscription")
	sub, err := (*a.stan).Subscribe(a.cfg.Stan.Channel, a.messageHandler)
	if err != nil {
		return err
	}
	defer sub.Close()
	log.Println("Stan channel subscription:\tSUCCESSFUL")

	log.Println("Start listening...")
	err = a.server.ListenAndServe()
	if err != nil {
		return err
	}
	return err
}

// Handle json-formated order from channel
// and store it in a cache and database
func (a *App) messageHandler(msg *stan.Msg) {
	// Add to cache
	data, err := model.NewFromByte(msg.Data)
	if err != nil {
		log.Println("Can't cache message: ", err)
		return
	}
	a.ch.Add(*data)

	err = a.db.InsertOrder(*data)
	if err != nil {
		log.Println("Can't insert data into database: ", err)
		return
	}
}
