package handlers

import "wedding-pictures/services"

type Handler struct {
	as services.AuthService
}

func NewHandler(as services.AuthService) *Handler {
	return &Handler{
		as: as,
	}
}
