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
	KafkaAddrs    []KafkaAddr `yaml:"kafka_addrs"`
	ConsumerGroup string      `yaml:"consumer_group"`
}

type KafkaAddr struct {
	Host  string `yaml:"host"`
	Port  string `yaml:"port"`
	Topic string `yaml:"topic"`
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
