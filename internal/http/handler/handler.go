package handler

import (
	"github.com/SicParv1sMagna/AtomHackMarsService/internal/http/repository"
)

type Handler struct {
	r *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{r: repo}
}
