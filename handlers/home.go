package handlers

import (
	"net/http"
	"wedding-pictures/views/auth"
	"wedding-pictures/views/home"
)

func HandleHome(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, home.Index())
}
func HandleHome2(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, auth.Login())
}
