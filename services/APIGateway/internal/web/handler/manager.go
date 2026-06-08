package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/domain"
	log "github.com/sirupsen/logrus"
)

type ManagerService interface {
	AddManager(ctx context.Context, manager domain.Manager) (uuid.UUID, error)
	RemoveManager(ctx context.Context, managerID uuid.UUID) error
}

type ManagerHandler struct {
	managerService ManagerService
}

func NewManagerHandler(managerService ManagerService) *ManagerHandler {
	return &ManagerHandler{managerService: managerService}
}

func (handler *ManagerHandler) AddManagerInfoHandler(w http.ResponseWriter, r *http.Request) {
	var newManager domain.Manager
	if err := json.NewDecoder(r.Body).Decode(&newManager); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		log.Error(err)
		return
	}

	employeeID, err := handler.managerService.AddManager(r.Context(), newManager)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error(err)
		return
	}

	type ManagerInfoResponse struct {
		ManagerID string `json:"managerId"`
	}

	if err = json.NewEncoder(w).Encode(ManagerInfoResponse{employeeID.String()}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error(err)
		return
	}
}
