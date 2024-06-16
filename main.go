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
		ImgSavePath:    "upload",
		MaxUploadBytes: 1024 * 1024 * 10, // 10 MB
	}

	as := services.NewAuthService()
	is := services.NewImageService(config)
	h := handlers.NewHandler(*as, *is, config)

	r := chi.NewMux()
	r.Use(middleware.RedirectSlashes)

	r.Handle("/*", public())

	// Serve/download gallery images
	r.Handle("/upload/*", http.StripPrefix("/upload/", http.FileServerFS(os.DirFS(config.ImgSavePath))))
	// r.Handle("/download/*", http.StripPrefix("/download/", http.FileServerFS(os.DirFS(config.ImgSavePath))))
	r.Get("/download/*", h.HandleDownload)

	// Pages
	r.Get("/", handlers.Make(h.HandleHome))
	r.Post("/upload", handlers.Make(h.HandleUpload))

	// Auth
	r.Get("/auth/{provider}", handlers.Make(h.HandleProviderLogin))
	r.Get("/auth/{provider}/callback", handlers.Make(h.HandleProviderCallback))
	r.Get("/auth/logout", handlers.Make(h.HandleProviderLogout))

	slog.Info("HTTP server started", "config", config)
	err := http.ListenAndServe(config.ListenAddr, r)
	if err != nil {
		slog.Error("failed starting server", "err", err)
		panic("cannot start server")
	}
}
