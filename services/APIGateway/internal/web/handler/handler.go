package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/domain"
)

type ApiService interface {
	Auth(ctx context.Context, userId uuid.UUID) (string, error)

	AddManager(ctx context.Context, manager domain.Manager) (uuid.UUID, error)
	RemoveManager(ctx context.Context, managerID uuid.UUID) error

	AddEmployee(ctx context.Context, employee domain.Employee) (uuid.UUID, error)
	GetEmployeesByManagerID(ctx context.Context, managerID uuid.UUID) ([]domain.Employee, error)
	RemoveEmployee(ctx context.Context, employeeID uuid.UUID) error

	AddPhoto(ctx context.Context, photo domain.Photo) (uuid.UUID, error)

	AddScreenshotStatistic(ctx context.Context, screenshot domain.ScreenshotStatistic) (uuid.UUID, error)
	GetScreenshotsStatistic(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]domain.ScreenshotStatistic, error)

	StartWorkSession(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error)
	StopWorkSession(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error)
	GetWorkSessions(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]domain.WorkSession, error)

	SaveToken(ctx context.Context, tokenID, userID string, expiresAt time.Time) error
	ExistsToken(ctx context.Context, tokenID string) (bool, error)
	DeleteToken(ctx context.Context, tokenID string) error
}

type ApiHandler struct {
	apiService ApiService
	jwtSecret  []byte
}

func NewApiHandler(apiService ApiService, jwtSecret []byte) *ApiHandler {
	return &ApiHandler{apiService: apiService, jwtSecret: jwtSecret}
}

func (handler *ApiHandler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	type TestRequest struct {
		UserID string `json:"userID"`
	}

	var req TestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)
		return
	}

	userId, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid manager uuid", http.StatusBadRequest)
		log.Error(err)
	}

	role, err := handler.apiService.Auth(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Error(err)
		return
	}

	accessToken, refreshToken, err := handler.generateTokens(userId.String(), role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Error(err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	type AuthResponse struct {
		UserID      string `json:"userID"`
		Role        string `json:"role"`
		AccessToken string `json:"accessToken"`
	}

	authResponse := AuthResponse{
		UserID:      userId.String(),
		Role:        role,
		AccessToken: accessToken,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(authResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
}

func (handler *ApiHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := extractRefreshToken(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)
		return
	}
	token, err := handler.parseRefreshToken(refreshToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)
		return
	}

	role, err := handler.apiService.Auth(r.Context(), uuid.MustParse(token.UserID))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Error(err)
		return
	}

	accessToken, refreshToken, err := handler.generateTokens(token.UserID, role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Error(err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	type AuthResponse struct {
		UserID      string `json:"userID"`
		Role        string `json:"role"`
		AccessToken string `json:"accessToken"`
	}

	authResponse := AuthResponse{
		UserID:      token.UserID,
		Role:        role,
		AccessToken: accessToken,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(authResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
}

func (handler *ApiHandler) AddScreenshotHandler(w http.ResponseWriter, r *http.Request) {
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

func (handler *ApiHandler) GetScreenshotsHandler(w http.ResponseWriter, r *http.Request) {
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

func (handler *ApiHandler) AddScreenshotStatisticHandler(w http.ResponseWriter, r *http.Request) {
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
	screenshotID, err := handler.apiService.AddScreenshotStatistic(r.Context(), screenshotStatistic)
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

func (handler *ApiHandler) GetScreenshotsStatisticHandler(w http.ResponseWriter, r *http.Request) {
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

	screenshots, err := handler.apiService.GetScreenshotsStatistic(r.Context(), employeeID, date)
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

func (handler *ApiHandler) StartWorkSessionHandler(w http.ResponseWriter, r *http.Request) {
	strEmployeeID := chi.URLParam(r, "employeeId")
	employeeID, err := uuid.Parse(strEmployeeID)
	if err != nil {
		http.Error(w, "invalid employee uuid", http.StatusBadRequest)
		log.Error(err)
		return
	}

	sessionID, err := handler.apiService.StartWorkSession(r.Context(), employeeID)
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

func (handler *ApiHandler) StopWorkSessionHandler(w http.ResponseWriter, r *http.Request) {
	strEmployeeID := chi.URLParam(r, "employeeId")
	employeeID, err := uuid.Parse(strEmployeeID)
	if err != nil {
		http.Error(w, "invalid employee uuid", http.StatusBadRequest)
		log.Error(err)
		return
	}

	sessionID, err := handler.apiService.StopWorkSession(r.Context(), employeeID)
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

func (handler *ApiHandler) GetWorkSessionsHandler(w http.ResponseWriter, r *http.Request) {
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

	sessions, err := handler.apiService.GetWorkSessions(r.Context(), employeeID, date)
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

func (handler *ApiHandler) AddEmployeeInfoHandler(w http.ResponseWriter, r *http.Request) {
	var newEmployee domain.Employee
	if err := json.NewDecoder(r.Body).Decode(&newEmployee); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		log.Error(err)
		return
	}

	employeeID, err := handler.apiService.AddEmployee(r.Context(), newEmployee)
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

func (handler *ApiHandler) AddManagerInfoHandler(w http.ResponseWriter, r *http.Request) {
	var newManager domain.Manager
	if err := json.NewDecoder(r.Body).Decode(&newManager); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		log.Error(err)
		return
	}

	employeeID, err := handler.apiService.AddManager(r.Context(), newManager)
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

func (handler *ApiHandler) AddEmployeePhoto(w http.ResponseWriter, r *http.Request) {
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

func (handler *ApiHandler) GetAllEmployeesInfoByManagerIDHandler(w http.ResponseWriter, r *http.Request) {
	strManagerID := chi.URLParam(r, "managerId")
	managerID, err := uuid.Parse(strManagerID)
	if err != nil {
		http.Error(w, "invalid manager uuid", http.StatusBadRequest)
		log.Error(err)
	}

	employees, err := handler.apiService.GetEmployeesByManagerID(r.Context(), managerID)
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

func (handler *ApiHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	strEmployeeID := chi.URLParam(r, "employeeId")
	employeeID, err := uuid.Parse(strEmployeeID)
	if err != nil {
		http.Error(w, "invalid employee uuid", http.StatusBadRequest)
		log.Error(err)
		return
	}

	if err = handler.apiService.RemoveEmployee(r.Context(), employeeID); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error(err)
		return
	}
}
