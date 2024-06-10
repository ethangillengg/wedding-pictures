package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"wedding-pictures/handlers"

	"wedding-pictures/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	as := services.NewAuthService()
	h := handlers.NewHandler(*as)

	r := chi.NewMux()
	r.Use(middleware.RedirectSlashes)
	r.Handle("/*", public())

	// Pages
	r.Get("/", handlers.Make(h.HandleHome))

	// Auth
	r.Get("/auth/{provider}", handlers.Make(h.HandleProviderLogin))
	r.Get("/auth/{provider}/callback", handlers.Make(h.HandleProviderCallback))
	r.Get("/auth/logout", handlers.Make(h.HandleProviderLogout))

	listenAddr := os.Getenv("LISTEN_ADDR")
	slog.Info("HTTP server started", "listenAddr", listenAddr)

	err := http.ListenAndServe(listenAddr, r)
	if err != nil {
		panic("cannot start server")
	}
}
