package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"wedding-pictures/views/home"

	"github.com/markbates/goth"
)

func (h *Handler) HandleHome(w http.ResponseWriter, r *http.Request) error {
	u, err := h.as.GetUserSession(r)
	if err != nil {
		u = goth.User{}
	}

	imgDirEntries, err := os.ReadDir(h.is.BaseDir)

	var imgs []string
	for _, entry := range imgDirEntries {
		if entry.Type().IsRegular() {
			imgs = append(imgs, filepath.Join(h.is.BaseDir, entry.Name()))
		}
	}

	return Render(w, r, home.Index(u, imgs))
}
