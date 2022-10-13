package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type (
	Config struct {
		Stan Stan
	}

	Stan struct {
		Clusterid string
		Userid    string
		Host      string
		Channel   string
	}
)

// Parse TOML file and return Config struct
func ParseConfig(filename string) (Config, error) {
	var config Config
	var text, err = os.ReadFile(filename)
	if err != nil {
		return config, err
	}
	_, err = toml.Decode(string(text[:]), &config)
	return config, err
}
