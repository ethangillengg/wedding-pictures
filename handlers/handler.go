package handlers

import (
	"wedding-pictures/services"
	"wedding-pictures/types"
)

type Handler struct {
	as  services.AuthService
	is  services.ImageService
	cfg types.Config
}

func NewHandler(as services.AuthService, is services.ImageService, cfg types.Config) *Handler {
	return &Handler{
		as:  as,
		is:  is,
		cfg: cfg,
	}
}
