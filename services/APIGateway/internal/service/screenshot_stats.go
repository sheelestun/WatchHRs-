package service

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/domain"
)

type ScreenshotStatisticStorage interface {
	AddScreenshotStatistic(ctx context.Context, screenshot domain.ScreenshotStatistic) (uuid.UUID, error)
	GetScreenshotsStatistic(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]domain.ScreenshotStatistic, error)
}

type screenshotStatisticService struct {
	screenshotStatisticStorage ScreenshotStatisticStorage
	validate                   *validator.Validate
}

func NewScreenshotStatisticService(ScreenshotStatisticStorage ScreenshotStatisticStorage, validate *validator.Validate) ScreenshotStatisticService {
	return &screenshotStatisticService{screenshotStatisticStorage: ScreenshotStatisticStorage, validate: validate}
}

func (s *screenshotStatisticService) AddScreenshotStatistic(ctx context.Context, screenshot domain.ScreenshotStatistic) (uuid.UUID, error) {
	screenshot.ID = uuid.New()
	screenshot.CreatedAt = time.Now()
	if err := s.validate.Struct(screenshot); err != nil {
		return uuid.Nil, err
	}
	return s.screenshotStatisticStorage.AddScreenshotStatistic(ctx, screenshot)
}

func (s *screenshotStatisticService) GetScreenshotsStatistic(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]domain.ScreenshotStatistic, error) {
	return s.screenshotStatisticStorage.GetScreenshotsStatistic(ctx, employeeID, date)
}
