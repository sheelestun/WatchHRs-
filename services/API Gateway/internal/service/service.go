package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/entity"
)

type Storage interface {
	AddManager(manager entity.Manager) (uuid.UUID, error)
	RemoveManager(managerId uuid.UUID) error

	AddEmployee(employee entity.Employee) (uuid.UUID, error)
	RemoveEmployee(employeeId uuid.UUID) error

	AddPhoto(photo entity.Photo) (uuid.UUID, error)

	AddScreenshotStatistic(screenshot entity.ScreenshotStatistic) (uuid.UUID, error)
	GetScreenshotsStatistic(employeeId uuid.UUID, date time.Time) ([]entity.ScreenshotStatistic, error)

	StartWorkSession(employeeId uuid.UUID) (uuid.UUID, error)
	StopWorkSession(employeeId uuid.UUID) (uuid.UUID, error)
	GetWorkSessions(employeeId uuid.UUID, date time.Time) ([]entity.WorkSession, error)
}

type ApiServiceImpl struct {
	storage Storage
}

func NewEmployeeService(storage Storage) *ApiServiceImpl {
	return &ApiServiceImpl{storage: storage}
}

func (a *ApiServiceImpl) AddManager(ctx context.Context, manager entity.Manager) (uuid.UUID, error) {
	return a.storage.AddManager(manager)
}

func (a *ApiServiceImpl) RemoveManager(ctx context.Context, managerId uuid.UUID) error {
	return a.storage.RemoveManager(managerId)
}

func (a *ApiServiceImpl) AddEmployee(ctx context.Context, employee entity.Employee) (uuid.UUID, error) {
	return a.storage.AddEmployee(employee)
}

func (a *ApiServiceImpl) RemoveEmployee(ctx context.Context, employeeId uuid.UUID) error {
	return a.storage.RemoveEmployee(employeeId)
}

func (a *ApiServiceImpl) AddPhoto(ctx context.Context, photo entity.Photo) (uuid.UUID, error) {
	return a.storage.AddPhoto(photo)
}

func (a *ApiServiceImpl) AddScreenshotStatistic(ctx context.Context, screenshot entity.ScreenshotStatistic) (uuid.UUID, error) {
	return a.storage.AddScreenshotStatistic(screenshot)
}

func (a *ApiServiceImpl) GetScreenshotsStatistic(ctx context.Context, employeeId uuid.UUID, date time.Time) ([]entity.ScreenshotStatistic, error) {
	return a.storage.GetScreenshotsStatistic(employeeId, date)
}

func (a *ApiServiceImpl) StartWorkSession(ctx context.Context, employeeId uuid.UUID) (uuid.UUID, error) {
	return a.storage.StartWorkSession(employeeId)
}

func (a *ApiServiceImpl) StopWorkSession(ctx context.Context, employeeId uuid.UUID) (uuid.UUID, error) {
	return a.storage.StopWorkSession(employeeId)
}

func (a *ApiServiceImpl) GetWorkSessions(ctx context.Context, employeeId uuid.UUID, date time.Time) ([]entity.WorkSession, error) {
	return a.storage.GetWorkSessions(employeeId, date)
}
