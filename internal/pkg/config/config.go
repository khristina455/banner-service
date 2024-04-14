package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServerConfig `yaml:"http_server"`
	PostgresConfig   `yaml:"postgres"`
	RedisConfig      `yaml:"redis"`
}

type HTTPServerConfig struct {
	Address           string        `yaml:"address" yaml-default:"localhost:8080"`
	Timeout           time.Duration `yaml:"timeout" yaml-default:"4s"`
	IDleTimeout       time.Duration `yaml:"idleTimeout" yaml-default:"60s"`
	ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout" yaml-defualt:"10s"`
	JWTSecret         string        `yaml:"JWTSecret"`
}

type RedisConfig struct {
	RedisAddr     string `yaml:"address"`
	RedisPassword string `yaml:"cache_pas"`
	DB            int
}

type PostgresConfig struct {
	DBName string `yaml:"db_name"`
	DBPass string `yaml:"db_pass"`
	DBHost string `yaml:"db_host"`
	DBPort int    `yaml:"db_port"`
	DBUser string `yaml:"db_user"`
}

func Load(filename string) (*Config, error) {
	var cfg Config

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Printf("config file does not exist: %s", filename)
		return nil, err
	}

	if err := cleanenv.ReadConfig(filename, &cfg); err != nil {
		log.Printf("cannot read %s: %v", filename, err)
		return nil, err
	}

	return &cfg, nil
}
