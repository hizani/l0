package main

import (
	"errors"
	"log"
	"os"
	"wbintern/l0/service/app"
	"wbintern/l0/service/config"
)

const usage = `Usage: service [FILE]`

func main() {
	// Get config filename as an argument
	cfgpath := "config.toml"

	if len(os.Args) > 1 {
		cfgpath = os.Args[1]
	}

	if _, err := os.Stat(cfgpath); errors.Is(err, os.ErrNotExist) {
		log.Fatalln(err, "\n", usage)
	}

	// Parse config
	var cfg, err = config.ParseConfig(cfgpath)
	if err != nil {
		log.Fatalln(err)
	}

	app, err := app.New(*cfg)
	if err != nil {
		log.Fatalln(err)
	}

	app.Run()

}
