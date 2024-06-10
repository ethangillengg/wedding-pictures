package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"wedding-pictures/services"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
)

func (h *Handler) HandleProviderLogin(w http.ResponseWriter, r *http.Request) error {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	if u, err := gothic.CompleteUserAuth(w, r); err == nil {
		h.as.StoreUserSession(w, r, u)
		http.Redirect(w, r, "/", http.StatusFound)
		return err
	}

	gothic.BeginAuthHandler(w, r)
	return nil
}

func (h *Handler) HandleProviderCallback(w http.ResponseWriter, r *http.Request) error {
	u, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return err
	}

	h.as.StoreUserSession(w, r, u)

	slog.Info("User logged in", "user", u.Email)
	http.Redirect(w, r, "/", http.StatusFound)

	return nil
}

func (h *Handler) HandleProviderLogout(w http.ResponseWriter, r *http.Request) error {
	h.as.ClearUserSession(w)
	gothic.Logout(w, r)
	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

func RequireAuth(h HTTPHandler, as *services.AuthService) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		session, err := as.GetUserSession(r)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return err
		}

		slog.Info("user is authenticated!", "email", session.Email)
		return h(w, r)
	}
}
