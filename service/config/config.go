package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type (
	Config struct {
		Stan     Stan
		Database Database
		Server   Server
	}

	Stan struct {
		Clusterid string
		Userid    string
		Host      string
		Port      int
		Channel   string
	}

	Database struct {
		Host string
		Port int
		User string
		Pass string
		Db   string
	}

	Server struct {
		Host string
		Port int
	}
)

// Parse TOML file and return Config struct
func ParseConfig(filename string) (*Config, error) {
	var config Config
	var text, err = os.ReadFile(filename)
	if err != nil {
		return &config, err
	}
	err = toml.Unmarshal(text, &config)
	return &config, err
}
