package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/domain"
	log "github.com/sirupsen/logrus"
)

type ScreenshotStatisticService interface {
	AddScreenshotStatistic(ctx context.Context, screenshot domain.ScreenshotStatistic) (uuid.UUID, error)
	GetScreenshotsStatistic(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]domain.ScreenshotStatistic, error)
}

type ScreenshotStatisticHandler struct {
	screenshotStatisticService ScreenshotStatisticService
}

func NewScreenshotStatisticHandler(ScreenshotStatisticService ScreenshotStatisticService) *ScreenshotStatisticHandler {
	return &ScreenshotStatisticHandler{screenshotStatisticService: ScreenshotStatisticService}
}

func (handler *ScreenshotStatisticHandler) AddScreenshotStatisticHandler(w http.ResponseWriter, r *http.Request) {
	strEmployeeID := chi.URLParam(r, "employeeId")
	employeeID, err := uuid.Parse(strEmployeeID)
	if err != nil {
		http.Error(w, "invalid employee uuid", http.StatusBadRequest)
		log.Error(err)
		return
	}

	var screenshotStatistic domain.ScreenshotStatistic
	if err = json.NewDecoder(r.Body).Decode(&screenshotStatistic); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		log.Error(err)
		return
	}

	screenshotStatistic.EmployeeID = employeeID
	screenshotID, err := handler.screenshotStatisticService.AddScreenshotStatistic(r.Context(), screenshotStatistic)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error(err)
		return
	}

	type StatisticResponse struct {
		ScreenshotID string `json:"screenshotId"`
	}
	statisticResponse := StatisticResponse{screenshotID.String()}

	if err := json.NewEncoder(w).Encode(statisticResponse); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error(err)
	}
}

func (handler *ScreenshotStatisticHandler) GetScreenshotsStatisticHandler(w http.ResponseWriter, r *http.Request) {
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

	screenshots, err := handler.screenshotStatisticService.GetScreenshotsStatistic(r.Context(), employeeID, date)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error(err)
		return
	}

	type StatisticResponse struct {
		Screenshots []domain.ScreenshotStatistic `json:"screenshots"`
	}
	statisticResponse := StatisticResponse{screenshots}
	if err := json.NewEncoder(w).Encode(statisticResponse); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error(err)
		return
	}
}
