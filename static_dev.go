// +build dev
//go:build dev
// +build dev

package main

import (
	"net/http"
	"os"
)

func public() http.Handler {
	return http.StripPrefix("/public/", http.FileServerFS(os.DirFS("public")))
}
