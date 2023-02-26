package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DB   *DB    `envconfig:"db"`
	Port string `envconfig:"app_port"`
}

type DB struct {
	Driver   string `envconfig:"driver" required:"true"`
	Password string `envconfig:"password"`
	User     string `envconfig:"user"`
	Name     string `envconfig:"name"`
	SSLMode  string `envconfig:"sslmode"`
}

func (c *Config) Process() error {
	return envconfig.Process("example", c)
}
