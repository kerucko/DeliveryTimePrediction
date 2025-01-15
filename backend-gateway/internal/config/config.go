package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
	Kafka    KafkaConfig    `yaml:"kafka"`
	Server   ServerConfig   `yaml:"server"`
}

type PostgresConfig struct {
	Host     string        `yaml:"host"`
	Port     string        `yaml:"port"`
	DBName   string        `yaml:"dbname"`
	User     string        `yaml:"user"`
	Password string        `yaml:"password"`
	Timeout  time.Duration `yaml:"timeout"`
}

type KafkaConfig struct {
	Brokers       []string     `yaml:"brokers"`
	Topics        TopicsConfig `yaml:"topics"`
	ConsumerGroup string       `yaml:"consumer_group"`
}

type TopicsConfig struct {
	Completed string `yaml:"completed"`
	Tasks     string `yaml:"tasks"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

func MustLoad() Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}
	var config Config
	err := cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		log.Fatalf("config not read: %v", err)
	}
	return config
}
