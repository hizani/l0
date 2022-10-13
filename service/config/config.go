package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type (
	Config struct {
		Stan     Stan
		Database Database
	}

	Stan struct {
		Clusterid string
		Userid    string
		Host      string
		Channel   string
	}

	Database struct {
		Host   string
		Port   string
		User   string
		Pass   string
		Db     string
		Driver string
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
