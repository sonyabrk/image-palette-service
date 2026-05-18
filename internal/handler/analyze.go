package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/sonyabrk/image-palette-service/internal/processor"
	"github.com/sonyabrk/image-palette-service/internal/worker"
)

const maxFileSize = 10 << 20

type errorResponse struct {
	Error string `json:"data"`
}

type successResponse struct {
	Data *processor.Result `json:"data"`
}

func Analyze(pool *worker.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)

		if err := r.ParseMultipartForm(32 << 20); err != nil {
			writeError(w, http.StatusBadRequest, "не удалось прочитать форму")
			return
		}

		file, header, err := r.FormFile("image")
		if err != nil {
			writeError(w, http.StatusBadRequest, "поле image обязательно")
			return
		}

		defer file.Close()

		contentType := header.Header.Get("Content-Type")
		if !isAllowedImageType(contentType) {
			writeError(w, http.StatusBadRequest,
				fmt.Sprintf("неподдерживаемый тип файла: %s", contentType))
			return
		}

		data, err := io.ReadAll(file)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "не удалось прочитать файл")
			return
		}

		result, err := pool.Submit(data)
		TrackRequest(err == nil)

		if err != nil {
			if errors.Is(err, ErrUnsupportedFormat) {
				writeError(w, http.StatusBadRequest, "неподдерживаемый формат изображения")
				return
			}

			writeError(w, http.StatusInternalServerError, "ошибка анализа изображения")
			return
		}

		writeJSON(w, http.StatusOK, successResponse{Data: result})
	}
}

func isAllowedImageType(contentType string) bool {
	allowed := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	return allowed[contentType]
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

var ErrUnsupportedFormat = errors.New("unsupported image format")
