package handlers

import (
	"net/http"
	"path/filepath"
	"wedding-pictures/views/components"
)

func (h *Handler) HandleUpload(w http.ResponseWriter, r *http.Request) error {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	// Get headers for filename, size and headers
	file, headers, err := r.FormFile("file")
	if err != nil {
		return err
	}

	// Upload the file
	if err := h.is.Upload(file, headers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return Render(w, r, components.GalleryImg(filepath.Join(h.is.FullDir, headers.Filename)))
}
