package server

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const maxFileSize = 10 << 20 // 10 MB

type UploadImageResponse struct {
	URL string `json:"url"`
}

func (h *Handler) UploadImage(r *http.Request) Response {
	log.Println("start uploading")
	err := r.ParseMultipartForm(maxFileSize) // 10 MB
	if err != nil {
		return JSONError(http.StatusBadRequest, "invalid multipart form, or file size exceeds limit")
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		return JSONError(http.StatusBadRequest, "image field is required")
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))

	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
	default:
		return JSONError(http.StatusBadRequest, "unsupported file type")
	}

	filename := randomFilename() + ext
	dstPath := filepath.Join(assetsDir+"/uploads", filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return JSONError(http.StatusInternalServerError, "failed to create file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return JSONError(http.StatusInternalServerError, "failed to save file")
	}
	url := "/" + assetsDir + "/uploads/" + filename
	return JSON(http.StatusCreated, UploadImageResponse{URL: url})
}

func randomFilename() string {
	b := make([]byte, 16)

	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}
