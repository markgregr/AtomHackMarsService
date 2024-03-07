package handler

import (
	"github.com/SicParv1sMagna/AtomHackMarsService/internal/http/repository"
	"github.com/SicParv1sMagna/AtomHackMarsService/internal/kafka"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	r *repository.Repository
	p *kafka.Producer
	c *kafka.Consumer
}

func New(repo *repository.Repository) *Handler {
	producer, err := kafka.NewProducer([]string{"localhost:29092"})
	if err != nil {
		log.Fatal("Error occured while creating producer", err)
	}

	consumer, err := kafka.NewConsumer([]string{"localhost:29092"})
	if err != nil {
		log.Fatal("Error occured while creating consumer")
	}

	return &Handler{
		r: repo,
		p: producer,
		c: consumer,
	}
}
