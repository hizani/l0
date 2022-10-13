package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"wbintern/l0/service/config"
)

const usage = `Usage: service [FILE]`

func main() {
	// Get config filename as an argument
	cfgpath := "config.toml"

	if len(os.Args) > 1 {
		cfgpath = os.Args[2]
	}

	if _, err := os.Stat(cfgpath); errors.Is(err, os.ErrNotExist) {
		log.Fatal(err, "\n", usage)
	}

	// Parse config
	var cfg, err = config.ParseConfig(cfgpath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)
	// Subscribe to nats streaming channel

	// Establish DB connection

	// Initialize cache store

	// Handle messages
}
