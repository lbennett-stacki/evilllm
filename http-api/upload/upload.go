package upload

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func UploadsPath(fileName string) string {
	return filepath.Join("_UPLOADS_", fileName)
}

func GeneratedPath(fileName string) string {
	return filepath.Join("_GENERATED_", fileName)
}

func UploadFile(r *http.Request) (path string, err error) {
	log.Debug("UploadFile called")

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Error("Error parsing multipart", "err", err)
		return path, err
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Error("Error retrieving file", "err", err)
		return path, err
	}
	defer file.Close()

	dst, err := os.Create(UploadsPath(handler.Filename))
	if err != nil {
		log.Error("Error creating dest file", "err", err)
		return path, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Error("Error copying to dest file", "err", err)
		return path, err
	}

	log.Debug("File uploaded successfully", "filename", handler.Filename)
	return handler.Filename, nil
}
