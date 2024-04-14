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
	JWTTTL            time.Duration `yaml:"JWTTTL" yaml-defualt:"6h"`
}

type RedisConfig struct {
	RedisAddr     string        `yaml:"address"`
	RedisPassword string        `yaml:"cachePas"`
	RedisTTL      time.Duration `yaml:"cacheTTL"`
}

type PostgresConfig struct {
	DBName string `yaml:"dbName"`
	DBPass string `yaml:"dbPass"`
	DBHost string `yaml:"dbHost"`
	DBPort int    `yaml:"dbPort"`
	DBUser string `yaml:"dbUser"`
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
