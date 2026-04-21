package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/domain"
	"github.com/sheelestun/WatchHRs-/internal/web/handler"
)

type ImageStorage interface {
	AddImage(ctx context.Context, photo domain.Image) (uuid.UUID, error)
}

type _ handler.ImageService

type imageService struct {
	imageStorage ImageStorage
	validate     *validator.Validate
}

func NewImageService(photoStorage ImageStorage, validate *validator.Validate) handler.ImageService {
	return &imageService{imageStorage: photoStorage, validate: validate}
}

func (s *imageService) AddImage(ctx context.Context, image domain.Image) (uuid.UUID, error) {
	image.ID = image.UserID
	if err := s.validate.Struct(image); err != nil {
		return uuid.Nil, err
	}
	return s.imageStorage.AddImage(ctx, image)
}
