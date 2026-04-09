package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ImageHandler struct {
	uploadDir string
}

func NewImageHandler(uploadDir string) *ImageHandler {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}
	return &ImageHandler{uploadDir: uploadDir}
}

func (h *ImageHandler) ServeImage(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")

	// Prevent path traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		http.Error(w, "invalid filename", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(h.uploadDir, filename)

	// Verify it's within the upload directory
	absPath, err := filepath.Abs(fullPath)
	if err != nil || !strings.HasPrefix(absPath, filepath.Clean(h.uploadDir)) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=86400")
	http.ServeFile(w, r, absPath)
}

// SaveUpload saves an uploaded file and returns the filename
func (h *ImageHandler) SaveUpload(r *http.Request, fieldName string, maxSize int64) (string, error) {
	r.Body = http.MaxBytesReader(nil, r.Body, maxSize)

	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Validate content type
	contentType := header.Header.Get("Content-Type")
	if !isAllowedImageType(contentType) {
		return "", &InvalidImageError{ContentType: contentType}
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := uuid.New().String() + ext

	dst, err := os.Create(filepath.Join(h.uploadDir, filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return filename, nil
}

func (h *ImageHandler) DeleteFile(filename string) {
	if filename == "" {
		return
	}
	os.Remove(filepath.Join(h.uploadDir, filename))
}

type InvalidImageError struct {
	ContentType string
}

func (e *InvalidImageError) Error() string {
	return "invalid image type: " + e.ContentType
}

func isAllowedImageType(ct string) bool {
	allowed := []string{"image/jpeg", "image/png", "image/webp", "image/gif"}
	for _, a := range allowed {
		if ct == a {
			return true
		}
	}
	return false
}
