package handler

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
}

func NewEmployeeHandler(employeeService EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{employeeService: employeeService}
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
		log.Error(err)
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

func (handler *EmployeeHandler) GetAllEmployeesInfoByManagerIDHandler(w http.ResponseWriter, r *http.Request) {
	strManagerID := chi.URLParam(r, "managerId")
	managerID, err := uuid.Parse(strManagerID)
	if err != nil {
		http.Error(w, "invalid manager uuid", http.StatusBadRequest)
		log.Error(err)
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
}
