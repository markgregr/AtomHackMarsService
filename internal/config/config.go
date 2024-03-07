package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func NewConfig() (*App, error) {
	var cfg App
	err := godotenv.Load()
    if err != nil {
        log.Fatalf("unable to load .env file: %e", err)
    }
	if err := env.Parse(&cfg); err != nil {
		panic(err)
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
