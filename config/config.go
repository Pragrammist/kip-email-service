package config

import (
	"bytes"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Smtp     Smtp
	Logger   Logger
	Grpc     Grpc
	Rabbit   Rabbit
	Database Database
}

type Logger struct {
	Mode string
}

type Smtp struct {
	Host     string
	Port     int
	User     string
	Password string
}

type Database struct {
	Host     string
	User     string
	Password string
	DbName   string
	Port     int
	SslMode  string
}

type Grpc struct {
	Port int
}

type Rabbit struct {
	Host        string
	Port        int
	User        string
	Password    string
	QueueName   string
	ConsumePool int
}

func LoadConfigFromEnv() (*Config, error) {
	v := viper.New()
	cfgEnv := os.Getenv("CONFIG")
	var c Config

	if cfgEnv == "" {
		log.Fatal("Provide config env variable")
	}

	v.AutomaticEnv()
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewReader([]byte(cfgEnv)))

	if err != nil {
		log.Fatalf("Cannot read config %s", err.Error())
	}

	err = v.Unmarshal(&c)

	if err != nil {
		log.Fatal(err.Error())
	}

	return &c, nil
}
