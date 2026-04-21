package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/domain"
	"github.com/sheelestun/WatchHRs-/internal/web/handler"
)

type PhotoStorage interface {
	AddPhoto(ctx context.Context, photo domain.Image) (uuid.UUID, error)
}

type _ handler.ImageService

type imageService struct {
	photoStorage PhotoStorage
	validate     *validator.Validate
}

func NewPhotoService(photoStorage PhotoStorage, validate *validator.Validate) handler.ImageService {
	return &imageService{photoStorage: photoStorage, validate: validate}
}

func (s *imageService) AddImage(ctx context.Context, image domain.Image) (uuid.UUID, error) {
	image.ID = image.UserID
	if err := s.validate.Struct(image); err != nil {
		return uuid.Nil, err
	}
	return s.photoStorage.AddPhoto(ctx, image)
}
