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

	// Create file
	defer file.Close()
	path := filepath.Join(is.FullDir, headers.Filename)
	dst, err := os.Create(path)
	defer dst.Close()
	if err != nil {
		return err
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		return err
	}

	go is.resize(headers.Filename)
	go is.convert(headers.Filename)

	return nil
}

func (is *ImageService) convert(fileName string) (string, error) {
	opts := bimg.Options{
		Crop: true,
		Type: bimg.WEBP,
	}

	path := filepath.Join("./", is.FullDir, fileName)

	slog.Debug("Converting", "path", path)

	buffer, err := bimg.Read(path)
	if err != nil {
		slog.Error("error reading image file", "err", err)
		return "", err
	}

	newImage, err := bimg.NewImage(buffer).Process(opts)
	if err != nil {
		slog.Error("error resizing image buffer", "err", err)
		return "", err
	}

	// replace ext with webp
	newName := strings.Replace(fileName, filepath.Ext(fileName), ".webp", 1)
	newPath := filepath.Join(is.BaseDir, newName)
	bimg.Write(newPath, newImage)

	return newName, nil
}

func (is *ImageService) resize(fileName string) {
	opts := bimg.Options{
		Width:  30,
		Height: 30,
		Crop:   true,
		Type:   bimg.WEBP,
	}

	path := filepath.Join("./", is.FullDir, fileName)

	slog.Debug("Resizing", "path", path)

	buffer, err := bimg.Read(path)
	if err != nil {
		slog.Error("error reading image file", "err", err)
	}

	newImage, err := bimg.NewImage(buffer).Process(opts)
	if err != nil {
		slog.Error("error resizing image buffer", "err", err)
	}

	// replace ext with webp
	newName := strings.Replace(fileName, filepath.Ext(fileName), ".webp", 1)
	bimg.Write(filepath.Join(is.SmallDir, newName), newImage)

}
