package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/domain"
)

type PhotoStorage interface {
	AddPhoto(ctx context.Context, photo domain.Photo) (uuid.UUID, error)
}

type PhotoServiceImpl struct {
	photoStorage PhotoStorage
	validate     *validator.Validate
}

func NewPhotoService(photoStorage PhotoStorage, validate *validator.Validate) *PhotoServiceImpl {
	return &PhotoServiceImpl{photoStorage: photoStorage, validate: validate}
}

func (s *PhotoServiceImpl) AddPhoto(ctx context.Context, photo domain.Photo) (uuid.UUID, error) {
	photo.ID = photo.UserID
	if err := s.validate.Struct(photo); err != nil {
		return uuid.Nil, err
	}
	return s.photoStorage.AddPhoto(ctx, photo)
}
