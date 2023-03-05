package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	JWTKeyword string `envconfig:"jwt_keyword"`
	DB         *DB    `envconfig:"db"`
	Port       string `envconfig:"app_port"`
	Minio      Minio  `envconfig:"minio"`
}

type DB struct {
	Driver   string `envconfig:"driver" required:"true"`
	Password string `envconfig:"password"`
	User     string `envconfig:"user"`
	Name     string `envconfig:"name"`
	SSLMode  string `envconfig:"sslmode"`
}

type Minio struct {
	KeyID     string `envconfig:"key_id"`
	SecretKey string `envconfig:"secret_key"`
	Endpoint  string `envconfig:"endpoint"`
	Bucket    string `envconfig:"bucket"`
}

func (c *Config) Process() error {
	return envconfig.Process("example", c)
}
