package handlers

import (
	"net/http"
	"wedding-pictures/views/home"

	"github.com/markbates/goth"
)

func (h *Handler) HandleHome(w http.ResponseWriter, r *http.Request) error {
	u, err := h.as.GetUserSession(r)
	if err != nil {
		u = goth.User{}
	}

	return Render(w, r, home.Index(u))
}
