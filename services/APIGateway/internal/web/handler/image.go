package handler

import (
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/client/cvimage"
	log "github.com/sirupsen/logrus"
)

type ImageHandler struct {
	cvClient *cvimage.Client
}

func NewImageHandler(cvClient *cvimage.Client) *ImageHandler {
	return &ImageHandler{cvClient: cvClient}
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
		if closeErr := file.Close(); closeErr != nil {
			log.Warnf("Error closing file: %v", closeErr)
		}
	}()

	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	employeeID := chi.URLParam(r, "employeeId")
	if _, err := uuid.Parse(employeeID); err != nil {
		http.Error(w, "invalid employee uuid", http.StatusBadRequest)
		log.Error(err)
		return
	}

	filename := header.Filename
	if filename == "" {
		filename = userID + "-" + uuid.NewString() + ".png"
	}

	resp, err := handler.cvClient.UploadScreenshot(r.Context(), employeeID, userID, file, filename)
	if err != nil {
		http.Error(w, "external service unavailable", http.StatusBadGateway)
		log.Error(err)
		return
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Warnf("Error closing body: %v", closeErr)
		}
	}()

	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func (handler *ImageHandler) proxyCVResponse(w http.ResponseWriter, resp *http.Response, defaultContentType string) {
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Warnf("Error closing body: %v", closeErr)
		}
	}()

	if contentType := resp.Header.Get("Content-Type"); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	} else if defaultContentType != "" {
		w.Header().Set("Content-Type", defaultContentType)
	}
	if disposition := resp.Header.Get("Content-Disposition"); disposition != "" {
		w.Header().Set("Content-Disposition", disposition)
	}

	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Error(err)
	}
}

func (handler *ImageHandler) GetScreenshotsHandler(w http.ResponseWriter, r *http.Request) {
	employeeID := chi.URLParam(r, "employeeId")
	if _, err := uuid.Parse(employeeID); err != nil {
		http.Error(w, "invalid employee uuid", http.StatusBadRequest)
		log.Error(err)
		return
	}

	date := chi.URLParam(r, "date")
	if _, err := time.Parse("2006-01-02", date); err != nil {
		http.Error(w, "invalid date", http.StatusBadRequest)
		log.Error(err)
		return
	}

	resp, err := handler.cvClient.GetScreenshots(r.Context(), employeeID, date)
	if err != nil {
		http.Error(w, "external service unavailable", http.StatusBadGateway)
		log.Error(err)
		return
	}
	handler.proxyCVResponse(w, resp, "application/json")
}

func (handler *ImageHandler) GetScreenshotsArchiveHandler(w http.ResponseWriter, r *http.Request) {
	employeeID := chi.URLParam(r, "employeeId")
	if _, err := uuid.Parse(employeeID); err != nil {
		http.Error(w, "invalid employee uuid", http.StatusBadRequest)
		log.Error(err)
		return
	}

	date := chi.URLParam(r, "date")
	if _, err := time.Parse("2006-01-02", date); err != nil {
		http.Error(w, "invalid date", http.StatusBadRequest)
		log.Error(err)
		return
	}

	resp, err := handler.cvClient.GetScreenshotsArchive(r.Context(), employeeID, date)
	if err != nil {
		http.Error(w, "external service unavailable", http.StatusBadGateway)
		log.Error(err)
		return
	}
	handler.proxyCVResponse(w, resp, "application/zip")
}

func (handler *ImageHandler) GetScreenshotFileHandler(w http.ResponseWriter, r *http.Request) {
	employeeID := chi.URLParam(r, "employeeId")
	if _, err := uuid.Parse(employeeID); err != nil {
		http.Error(w, "invalid employee uuid", http.StatusBadRequest)
		log.Error(err)
		return
	}

	filename := chi.URLParam(r, "filename")
	resp, err := handler.cvClient.GetScreenshotFile(r.Context(), employeeID, filename)
	if err != nil {
		http.Error(w, "external service unavailable", http.StatusBadGateway)
		log.Error(err)
		return
	}
	handler.proxyCVResponse(w, resp, "image/png")
}
