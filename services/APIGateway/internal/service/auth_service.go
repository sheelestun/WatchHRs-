package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AuthStorage interface {
	FindUser(ctx context.Context, userId uuid.UUID) (string, error)
}

type AuthCache interface {
	SaveTokenInCache(ctx context.Context, tokenID, userID string, expiresAt time.Time) error
	ExistsTokenInCache(ctx context.Context, tokenID string) (bool, error)
	DeleteTokenInCache(ctx context.Context, tokenID string) error
}

type AuthServiceImpl struct {
	authStorage AuthStorage
	cache       AuthCache
}

func NewAuthService(authStorage AuthStorage, cache AuthCache) *AuthServiceImpl {
	return &AuthServiceImpl{
		authStorage: authStorage,
		cache:       cache,
	}
}

func (a *AuthServiceImpl) Auth(ctx context.Context, userId uuid.UUID) (string, error) {
	return a.authStorage.FindUser(ctx, userId)
}

func (a *AuthServiceImpl) DeleteToken(ctx context.Context, tokenID string) error {
	return a.cache.DeleteTokenInCache(ctx, tokenID)
}

func (a *AuthServiceImpl) ExistsToken(ctx context.Context, tokenID string) (bool, error) {
	return a.cache.ExistsTokenInCache(ctx, tokenID)
}

func (a *AuthServiceImpl) SaveToken(ctx context.Context, tokenID, userID string, expiresAt time.Time) error {
	return a.cache.SaveTokenInCache(ctx, tokenID, userID, expiresAt)
}
