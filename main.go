package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"wedding-pictures/handlers"
	"wedding-pictures/types"

	"wedding-pictures/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	config := types.Config{
		ListenAddr:     os.Getenv("LISTEN_ADDR"),
		ImgSavePath:    "./uploads",
		MaxUploadBytes: 1024 * 1024 * 10,
	}

	err := os.Mkdir(config.ImgSavePath, 0750)
	if os.IsExist(err) {
		slog.Info("upload dir already exists", "dir", config.ImgSavePath)
	} else if err != nil {
		log.Fatal(err)
	}

	as := services.NewAuthService()
	h := handlers.NewHandler(*as, config)

	r := chi.NewMux()
	r.Use(middleware.RedirectSlashes)

	r.Handle("/*", public())
	// Serve gallery images
	r.Handle("/gallery/*", http.StripPrefix("/gallery/", http.FileServerFS(os.DirFS(config.ImgSavePath))))

	// Pages
	r.Get("/", handlers.Make(h.HandleHome))
	r.Post("/upload", handlers.Make(h.HandleUpload))

	// Auth
	r.Get("/auth/{provider}", handlers.Make(h.HandleProviderLogin))
	r.Get("/auth/{provider}/callback", handlers.Make(h.HandleProviderCallback))
	r.Get("/auth/logout", handlers.Make(h.HandleProviderLogout))

	slog.Info("HTTP server started", "config", config)
	err = http.ListenAndServe(config.ListenAddr, r)
	if err != nil {
		slog.Error("failed starting server", "err", err)
		panic("cannot start server")
	}
}
