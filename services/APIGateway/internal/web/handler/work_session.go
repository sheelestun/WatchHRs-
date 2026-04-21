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

type WorkSessionService interface {
	StartWorkSession(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error)
	StopWorkSession(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error)
	GetWorkSessions(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]domain.WorkSession, error)
}

type WorkSessionHandler struct {
	workSessionService WorkSessionService
}

func NewWorkSessionHandler(workSessionService WorkSessionService) *WorkSessionHandler {
	return &WorkSessionHandler{workSessionService: workSessionService}
}

func (handler *WorkSessionHandler) StartWorkSessionHandler(w http.ResponseWriter, r *http.Request) {
	strEmployeeID := chi.URLParam(r, "employeeId")
	employeeID, err := uuid.Parse(strEmployeeID)
	if err != nil {
		http.Error(w, "invalid employee uuid", http.StatusBadRequest)
		log.Error(err)
		return
	}

	sessionID, err := handler.workSessionService.StartWorkSession(r.Context(), employeeID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error(err)
		return
	}

	type SessionResponse struct {
		SessionID string `json:"sessionId"`
	}

	sessionResponse := SessionResponse{sessionID.String()}
	if err := json.NewEncoder(w).Encode(sessionResponse); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error(err)
		return
	}
}

func (handler *WorkSessionHandler) StopWorkSessionHandler(w http.ResponseWriter, r *http.Request) {
	strEmployeeID := chi.URLParam(r, "employeeId")
	employeeID, err := uuid.Parse(strEmployeeID)
	if err != nil {
		http.Error(w, "invalid employee uuid", http.StatusBadRequest)
		log.Error(err)
		return
	}

	sessionID, err := handler.workSessionService.StopWorkSession(r.Context(), employeeID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error(err)
		return
	}

	type SessionResponse struct {
		SessionID string `json:"sessionId"`
	}
	sessionResponse := SessionResponse{sessionID.String()}
	if err := json.NewEncoder(w).Encode(sessionResponse); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error(err)
		return
	}
}

func (handler *WorkSessionHandler) GetWorkSessionsHandler(w http.ResponseWriter, r *http.Request) {
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

	sessions, err := handler.workSessionService.GetWorkSessions(r.Context(), employeeID, date)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error(err)
		return
	}

	type SessionResponse struct {
		WorkSessions []domain.WorkSession `json:"workSessions"`
	}
	sessionResponse := SessionResponse{sessions}
	if err := json.NewEncoder(w).Encode(sessionResponse); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error(err)
		return
	}
}
