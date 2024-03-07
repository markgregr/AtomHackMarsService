package config

import (
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func NewConfig() (*App, error) {
	var cfg App
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("unable to load .env file: %v", err)
		}
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("error parsing environment variables: %v", err)
	}
	log.Println("cfg: ", cfg)
	level, err := log.ParseLevel(cfg.ErrorLevel)
	if err != nil {
		panic(err)
	}

	log.SetLevel(level)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)

	return &cfg, nil
}
