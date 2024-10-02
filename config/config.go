package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App   App   `yaml:"app"`
	HTTP  HTTP  `yaml:"http"`
	DB    DB    `yaml:"postgres"`
	Kafka Kafka `yaml:"kafka"`
}
type App struct {
	Name    string `yaml:"name" env-required:"true" env:"APP_NAME"`
	Version string `yaml:"version" env-required:"true" env:"APP_VERSION"`
}
type HTTP struct {
	Host string `yaml:"host" env-requiered:"true" env:"HTTP_HOST"`
	Port string `yaml:"port" env-required:"true" env:"HTTP_PORT"`
}
type DB struct {
	Host     string `yaml:"host" env-requiered:"true" env:"PG_HOST"`
	Port     string `yaml:"port" env-required:"true" env:"PG_PORT"`
	User     string `yaml:"user" env-required:"true" env:"PG_USER"`
	Password string `yaml:"password" env-required:"true" env:"PG_PASSWORD"`
	DBName   string `yaml:"name" env:"PG_NAME" env-required:"true" `
	PgDriver string `yaml:"pg_driver" env:"PG_PG_DRIVER" env-required:"true" `
	Schema   string `yaml:"schema" env:"PG_SCHEMA" env-required:"true" `
}
type Kafka struct {
	BootstrapServers string `env-required:"true" yaml:"bootstrap_servers" env:"KAFKA_BOOTSTRAP_SERVERS"`
	Topic            string `yaml:"topic" env-required:"true" env:"KAFKA_TOPIC"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "./config/config.yml"
	}
	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("сonfiguration error: %v", err)
	}
	return cfg, nil
}
