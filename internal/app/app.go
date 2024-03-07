package app

import (
	"github.com/SicParv1sMagna/AtomHackMarsService/internal/config"
	"github.com/SicParv1sMagna/AtomHackMarsService/internal/http/handler"
	"github.com/SicParv1sMagna/AtomHackMarsService/internal/http/repository"
)

type Application struct {
	cfg     *config.App
	handler *handler.Handler
}

func New() (*Application, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	repo, err := repository.New(cfg)
	if err != nil {
		return nil, err
	}

	h := handler.New(repo, cfg)

	app := &Application{
		cfg:     cfg,
		handler: h,
	}

	return app, nil
}
