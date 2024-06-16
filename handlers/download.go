package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	// "time"
)

func (h *Handler) HandleDownload(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimSuffix(filepath.Base(r.URL.Path), filepath.Ext(r.URL.Path))
	glob := filepath.Join(h.is.FullDir, filename+".*")
	paths, err := filepath.Glob(glob)

	slog.Info("download", "path", r.URL.Path, "filename", filename, "glob", glob)

	if err != nil {
		slog.Error("error globbing", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if len(paths) == 0 {
		slog.Info("no globs")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	file, err := os.Open(paths[0])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(paths[0]))
	w.Header().Set("Content-Type", "application/octet-stream")

	io.Copy(w, file)
}
