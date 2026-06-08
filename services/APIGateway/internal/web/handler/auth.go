package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/client/cvimage"
	log "github.com/sirupsen/logrus"
)

type AuthService interface {
	Auth(ctx context.Context, userId uuid.UUID) (string, error)

	SaveToken(ctx context.Context, tokenID, userID string, expiresAt time.Time) error
	ExistsToken(ctx context.Context, tokenID string) (bool, error)
	DeleteToken(ctx context.Context, tokenID string) error
}

type AuthHandler struct {
	authService AuthService
	cvClient    *cvimage.Client
	jwtSecret   []byte
}

func NewAuthHandler(authService AuthService, cvClient *cvimage.Client, jwtSecret []byte) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cvClient:    cvClient,
		jwtSecret:   jwtSecret,
	}
}

func (handler *AuthHandler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	const maxFileSize = 15 << 20

	contentType := r.Header.Get("Content-Type")
	if contentType == "" || contentType == "application/json" {
		handler.authByUserID(w, r)
		return
	}
	if strings.HasPrefix(contentType, "application/json") {
		handler.authByUserID(w, r)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)
	if err := r.ParseMultipartForm(maxFileSize); err == nil {
		file, _, err := r.FormFile("photo")
		if err == nil {
			defer file.Close()
			handler.authByPhoto(w, r, file)
			return
		}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "photo is required", http.StatusBadRequest)
		return
	}
	handler.authByPhoto(w, r, bytes.NewReader(body))
}

func (handler *AuthHandler) authByUserID(w http.ResponseWriter, r *http.Request) {
	type AuthRequest struct {
		UserID string `json:"userID"`
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid user uuid", http.StatusBadRequest)
		log.Error(err)
		return
	}

	handler.issueTokens(w, r, userID)
}

func (handler *AuthHandler) authByPhoto(w http.ResponseWriter, r *http.Request, body io.Reader) {
	result, status, err := handler.cvClient.Authenticate(r.Context(), body)
	if err != nil {
		http.Error(w, "external service unavailable", http.StatusBadGateway)
		log.Error(err)
		return
	}
	if status != http.StatusOK || result == nil {
		http.Error(w, "authentication failed", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(result.UserID)
	if err != nil {
		http.Error(w, "invalid user id from cv service", http.StatusInternalServerError)
		log.Error(err)
		return
	}

	handler.issueTokens(w, r, userID)
}

func (handler *AuthHandler) issueTokens(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	role, err := handler.authService.Auth(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Error(err)
		return
	}

	accessToken, refreshToken, err := handler.generateTokens(userID.String(), role)
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
		UserID:      userID.String(),
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

func (handler *AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
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

	role, err := handler.authService.Auth(r.Context(), uuid.MustParse(token.UserID))
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
