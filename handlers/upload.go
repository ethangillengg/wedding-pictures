package handlers

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

func (h *Handler) HandleUpload(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Trying to upload file")
	// u, err := h.as.GetUserSession(r)
	// if err != nil {
	// 	u = goth.User{}
	// }

	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	// Get fileHandler for filename, size and headers
	file, fileHandler, err := r.FormFile("file")
	if err != nil {
		return err
	}

	defer file.Close()
	slog.Info("File Name", "FileName", fileHandler.Filename)
	slog.Info("File Size", "Size", fileHandler.Size)
	slog.Info("MIME Header", "Header", fileHandler.Header)

	// Create file
	imgFilePath := filepath.Join(h.cfg.ImgSavePath, fileHandler.Filename)
	dst, err := os.Create(imgFilePath)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	fmt.Fprintf(w, "Successfully Uploaded File!\n")

	return nil
	// return Render(w, r, home.Index(u))
}
