package services

import (
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"wedding-pictures/types"

	"github.com/h2non/bimg"
)

const (
	SmallDir = "small"
	FullDir  = "full"
)

type ImageService struct {
	SmallDir       string
	FullDir        string
	BaseDir        string
	MaxUploadBytes int
}

func NewImageService(cfg types.Config) *ImageService {
	base := cfg.ImgSavePath
	small := filepath.Join(cfg.ImgSavePath, SmallDir)
	full := filepath.Join(cfg.ImgSavePath, FullDir)

	err := os.Mkdir(base, 0750)
	if os.IsExist(err) {
		slog.Info("upload dir already exists", "dir", base)
	} else if err != nil {
		panic("failed to make upload dir")
	}

	err = os.Mkdir(small, 0750)
	if os.IsExist(err) {
		slog.Info("small dir already exists", "dir", small)
	} else if err != nil {
		panic("failed to make small dir")
	}
	err = os.Mkdir(full, 0750)
	if os.IsExist(err) {
		slog.Info("full dir already exists", "dir", full)
	} else if err != nil {
		panic("failed to make full dir")
	}

	return &ImageService{
		BaseDir:        base,
		SmallDir:       small,
		FullDir:        full,
		MaxUploadBytes: cfg.MaxUploadBytes,
	}
}

func (is *ImageService) Upload(file multipart.File, headers *multipart.FileHeader) error {
	slog.Info("Uploading file", "name", headers.Filename, "size", headers.Size)
	if headers.Size > int64(is.MaxUploadBytes) {
		return fmt.Errorf("file \"%v\" (%v bytes) exceeds the upload limit of %v bytes", headers.Filename, headers.Size, is.MaxUploadBytes)
	}

	name := headers.Filename

	fullPath := filepath.Join(is.FullDir, name)
	smallPath := filepath.Join(is.SmallDir, strings.Replace(name, filepath.Ext(name), ".webp", 1))
	basePath := filepath.Join(is.BaseDir, strings.Replace(name, filepath.Ext(name), ".webp", 1))

	// Create file
	dst, err := os.Create(fullPath)
	defer file.Close()
	defer dst.Close()
	if err != nil {
		return err
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		return err
	}
	dst.Sync()

	// save 30x30 blurred thumbnail
	go is.process(fullPath, smallPath, bimg.Options{
		Width:   10,
		Height:  10,
		Quality: 1,
		Crop:    true,
		Gravity: bimg.GravitySmart,
		Type:    bimg.WEBP,
	})
	//
	// save web-friendly version
	go is.process(fullPath, basePath, bimg.Options{
		Width:  1000,
		Height: 1000,

		Crop:    true,
		Gravity: bimg.GravitySmart,
		Type:    bimg.WEBP,
	})

	slog.Info("Saved full image to disk", "fullPath", fullPath)

	return nil
}

func (is *ImageService) process(src string, dest string, opts bimg.Options) error {
	slog.Info("Processing image", "dest", dest)

	bytes, err := bimg.Read(src)
	if err != nil {
		slog.Error("error reading image file", "err", err)
	}

	img, err := bimg.NewImage(bytes).Process(opts)
	if err != nil {
		slog.Error("error processing image", "dest", dest, "err", err)
		return err
	}

	bimg.Write(dest, img)
	return nil
}
