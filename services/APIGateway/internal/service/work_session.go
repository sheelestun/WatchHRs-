package service

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/domain"
)

type WorkSessionStorage interface {
	StartWorkSession(ctx context.Context, session domain.WorkSession) (uuid.UUID, error)
	StopWorkSession(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error)
	GetWorkSessions(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]domain.WorkSession, error)
}

type WorkSessionServiceImpl struct {
	workSessionStorage WorkSessionStorage
	validate           *validator.Validate
}

func NewWorkSessionServiceImpl(WorkSessionStorage WorkSessionStorage, validate *validator.Validate) *WorkSessionServiceImpl {
	return &WorkSessionServiceImpl{workSessionStorage: WorkSessionStorage, validate: validate}
}

func (w *WorkSessionServiceImpl) StartWorkSession(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error) {
	newSession := domain.WorkSession{ID: uuid.New(), EmployeeID: employeeID, StartTime: time.Now()}
	if err := w.validate.Struct(newSession); err != nil {
		return uuid.Nil, err
	}
	return w.workSessionStorage.StartWorkSession(ctx, newSession)
}

func (w *WorkSessionServiceImpl) StopWorkSession(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error) {
	return w.workSessionStorage.StopWorkSession(ctx, employeeID)
}

func (w *WorkSessionServiceImpl) GetWorkSessions(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]domain.WorkSession, error) {
	return w.workSessionStorage.GetWorkSessions(ctx, employeeID, date)
}
