package handlers

import (
	"wedding-pictures/services"
	"wedding-pictures/types"
)

type Handler struct {
	as  services.AuthService
	cfg types.Config
}

func NewHandler(as services.AuthService, cfg types.Config) *Handler {
	return &Handler{
		as:  as,
		cfg: cfg,
	}
}
