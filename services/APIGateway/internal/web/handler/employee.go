package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/client/cvimage"
	"github.com/sheelestun/WatchHRs-/internal/domain"
	log "github.com/sirupsen/logrus"
)

type EmployeeService interface {
	AddEmployee(ctx context.Context, employee domain.Employee) (uuid.UUID, error)
	GetEmployeesByManagerID(ctx context.Context, managerID uuid.UUID) ([]domain.Employee, error)
	RemoveEmployee(ctx context.Context, employeeID uuid.UUID) error
}

type EmployeeHandler struct {
	employeeService EmployeeService
	cvClient        *cvimage.Client
}

func NewEmployeeHandler(employeeService EmployeeService, cvClient *cvimage.Client) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: employeeService,
		cvClient:        cvClient,
	}
}

func (handler *EmployeeHandler) AddEmployeeInfoHandler(w http.ResponseWriter, r *http.Request) {
	var newEmployee domain.Employee
	if err := json.NewDecoder(r.Body).Decode(&newEmployee); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		log.Error(err)
		return
	}

	employeeID, err := handler.employeeService.AddEmployee(r.Context(), newEmployee)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error(err)
		return
	}

	type EmployeeInfoResponse struct {
		EmployeeID string `json:"employeeId"`
	}

	if err = json.NewEncoder(w).Encode(EmployeeInfoResponse{employeeID.String()}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error(err)
		return
	}
}

func (handler *EmployeeHandler) AddEmployeePhoto(w http.ResponseWriter, r *http.Request) {
	const maxFileSize = 15 << 20 // 15MB

	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)

	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		log.Error(err)
		return
	}

	userID := r.FormValue("userId")
	if userID == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}
	if _, err := uuid.Parse(userID); err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		log.Error(err)
		return
	}

	file, _, err := r.FormFile("screenshot")
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

	resp, err := handler.cvClient.UploadPhoto(r.Context(), userID, file)
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

func (handler *EmployeeHandler) GetAllEmployeesInfoByManagerIDHandler(w http.ResponseWriter, r *http.Request) {
	strManagerID := chi.URLParam(r, "managerId")
	managerID, err := uuid.Parse(strManagerID)
	if err != nil {
		http.Error(w, "invalid manager uuid", http.StatusBadRequest)
		log.Error(err)
		return
	}

	employees, err := handler.employeeService.GetEmployeesByManagerID(r.Context(), managerID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error(err)
		return
	}

	if err = json.NewEncoder(w).Encode(employees); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error(err)
		return
	}
}

func (handler *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	strEmployeeID := chi.URLParam(r, "employeeId")
	employeeID, err := uuid.Parse(strEmployeeID)
	if err != nil {
		http.Error(w, "invalid employee uuid", http.StatusBadRequest)
		log.Error(err)
		return
	}

	if err = handler.employeeService.RemoveEmployee(r.Context(), employeeID); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error(err)
		return
	}

	if err = handler.cvClient.DeletePhoto(r.Context(), employeeID.String()); err != nil {
		log.WithError(err).Warn("failed to delete employee photo from cv storage")
	}
}
