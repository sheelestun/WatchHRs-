package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
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
	jwtSecret   []byte
}

func NewAuthHandler(authService AuthService, jwtSecret []byte) *AuthHandler {
	return &AuthHandler{authService: authService, jwtSecret: jwtSecret}
}

func (handler *AuthHandler) AuthHandler(w http.ResponseWriter, r *http.Request) {
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

	role, err := handler.authService.Auth(r.Context(), userId)
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
