package handler

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/domain"
	log "github.com/sirupsen/logrus"
)

type ImageService interface {
	AddImage(ctx context.Context, photo domain.Image) (uuid.UUID, error)
}

type ImageHandler struct {
	imageService ImageService
}

func NewImageHandler(ImageService ImageService) *ImageHandler {
	return &ImageHandler{imageService: ImageService}
}

func (handler *ImageHandler) AddScreenshotHandler(w http.ResponseWriter, r *http.Request) {
	const maxFileSize = 15 << 20 // 15MB

	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)

	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		log.Error(err)
		return
	}

	file, header, err := r.FormFile("screenshot")
	if err != nil {
		http.Error(w, "screenshot is required", http.StatusBadRequest)
		log.Error(err)
		return
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Warnf("Error closing file: %v", err)
		}
	}()

	// Потоковая передача во внешний сервис
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer func() {
			err = pw.Close()
			if err != nil {
				log.Warnf("Error closing writer: %v", err)
			}
		}()
		defer func() {
			err = writer.Close()
			if err != nil {
				log.Warnf("Error closing writer: %v", err)
			}
		}()

		part, err := writer.CreateFormFile("screenshot", header.Filename)
		if err != nil {
			_ = pw.CloseWithError(err)
			return
		}

		if _, err := io.Copy(part, file); err != nil {
			_ = pw.CloseWithError(err)
			return
		}
	}()

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"localhost:8081/photo", // TODO: Вынести ссылку на сервис в конфиг
		pr,
	)
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		log.Error(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "external service unavailable", http.StatusBadGateway)
		return
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Warnf("Error closing body: %v", err)
		}
	}()

	w.WriteHeader(resp.StatusCode)
}

func (handler *ImageHandler) GetScreenshotsHandler(w http.ResponseWriter, r *http.Request) {
	strEmployeeID := chi.URLParam(r, "employeeId")
	employeeID, err := uuid.Parse(strEmployeeID)
	if err != nil {
		http.Error(w, "invalid employee uuid", http.StatusBadRequest)
		log.Error(err)
		return
	}

	strDate := chi.URLParam(r, "date")
	date, err := time.Parse("2006-01-02", strDate)
	if err != nil {
		http.Error(w, "invalid date", http.StatusBadRequest)
		log.Error(err)
		return
	}

	// Формируем URL внешнего сервиса
	externalURL := fmt.Sprintf(
		"localhost:8081/screenshots/%s/%s", // TODO: Вынести ссылку на сервис в конфиг
		employeeID.String(),
		date.Format("2006-01-02"),
	)

	resp, err := http.Get(externalURL)
	if err != nil {
		http.Error(w, "external service unavailable", http.StatusBadGateway)
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Warnf("Error closing body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "external service error", http.StatusBadGateway)
		return
	}

	// Проксируем ответ (массив скриншотов)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "external service unavailable", http.StatusBadGateway)
		log.Error(err)
		return
	}
}
