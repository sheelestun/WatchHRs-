package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/web/handler"
)

type AuthStorage interface {
	FindUser(ctx context.Context, userId uuid.UUID) (string, error)
}

type AuthCache interface {
	SaveTokenInCache(ctx context.Context, tokenID, userID string, expiresAt time.Time) error
	ExistsTokenInCache(ctx context.Context, tokenID string) (bool, error)
	DeleteTokenInCache(ctx context.Context, tokenID string) error
}

type _ handler.AuthService

type authService struct {
	authStorage AuthStorage
	cache       AuthCache
}

func NewAuthService(authStorage AuthStorage, cache AuthCache) handler.AuthService {
	return &authService{
		authStorage: authStorage,
		cache:       cache,
	}
}

func (a *authService) Auth(ctx context.Context, userId uuid.UUID) (string, error) {
	return a.authStorage.FindUser(ctx, userId)
}

func (a *authService) DeleteToken(ctx context.Context, tokenID string) error {
	return a.cache.DeleteTokenInCache(ctx, tokenID)
}

func (a *authService) ExistsToken(ctx context.Context, tokenID string) (bool, error) {
	return a.cache.ExistsTokenInCache(ctx, tokenID)
}

func (a *authService) SaveToken(ctx context.Context, tokenID, userID string, expiresAt time.Time) error {
	return a.cache.SaveTokenInCache(ctx, tokenID, userID, expiresAt)
}
