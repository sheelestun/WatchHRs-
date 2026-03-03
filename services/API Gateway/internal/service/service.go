package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/entity"
)

// Storage interface for repository
type Storage interface {
	FindUser(userId uuid.UUID) (string, error)

	AddManager(manager entity.Manager) (uuid.UUID, error)
	RemoveManager(managerID uuid.UUID) error

	AddEmployee(employee entity.Employee) (uuid.UUID, error)
	GetAllEmployeesByManagerID(managerID uuid.UUID) ([]entity.Employee, error)
	RemoveEmployee(employeeID uuid.UUID) error

	AddPhoto(photo entity.Photo) (uuid.UUID, error)

	AddScreenshotStatistic(screenshot entity.ScreenshotStatistic) (uuid.UUID, error)
	GetScreenshotsStatistic(employeeID uuid.UUID, date time.Time) ([]entity.ScreenshotStatistic, error)

	StartWorkSession(employeeID uuid.UUID) (uuid.UUID, error)
	StopWorkSession(employeeID uuid.UUID) (uuid.UUID, error)
	GetWorkSessions(employeeID uuid.UUID, date time.Time) ([]entity.WorkSession, error)

	SaveToken(tokenID, userID string, expiresAt time.Time) error
	ExistsToken(tokenID string) (bool, error)
	DeleteToken(tokenID string) error
}

// APIServiceImpl implementation for ApiService interface
type APIServiceImpl struct {
	storage Storage
}

func (a *APIServiceImpl) Auth(ctx context.Context, userId uuid.UUID) (string, error) {
	return a.storage.FindUser(userId)
}

func (a *APIServiceImpl) DeleteToken(tokenID string) error {
	return a.storage.DeleteToken(tokenID)
}

func (a *APIServiceImpl) ExistsToken(tokenID string) (bool, error) {
	return a.storage.ExistsToken(tokenID)
}

func (a *APIServiceImpl) SaveToken(tokenID, userID string, expiresAt time.Time) error {
	return a.storage.SaveToken(tokenID, userID, expiresAt)
}

// NewAPIService constructor
func NewAPIService(storage Storage) *APIServiceImpl {
	return &APIServiceImpl{storage: storage}
}

func (a *APIServiceImpl) AddManager(ctx context.Context, manager entity.Manager) (uuid.UUID, error) {
	return a.storage.AddManager(manager)
}

func (a *APIServiceImpl) RemoveManager(ctx context.Context, managerID uuid.UUID) error {
	return a.storage.RemoveManager(managerID)
}

func (a *APIServiceImpl) AddEmployee(ctx context.Context, employee entity.Employee) (uuid.UUID, error) {
	return a.storage.AddEmployee(employee)
}

func (a *APIServiceImpl) GetEmployeesByManagerID(managerID uuid.UUID) ([]entity.Employee, error) {
	return a.storage.GetAllEmployeesByManagerID(managerID)
}

func (a *APIServiceImpl) RemoveEmployee(ctx context.Context, employeeID uuid.UUID) error {
	return a.storage.RemoveEmployee(employeeID)
}

func (a *APIServiceImpl) AddPhoto(ctx context.Context, photo entity.Photo) (uuid.UUID, error) {
	return a.storage.AddPhoto(photo)
}

func (a *APIServiceImpl) AddScreenshotStatistic(ctx context.Context, screenshot entity.ScreenshotStatistic) (uuid.UUID, error) {
	return a.storage.AddScreenshotStatistic(screenshot)
}

func (a *APIServiceImpl) GetScreenshotsStatistic(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]entity.ScreenshotStatistic, error) {
	return a.storage.GetScreenshotsStatistic(employeeID, date)
}

func (a *APIServiceImpl) StartWorkSession(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error) {
	return a.storage.StartWorkSession(employeeID)
}

func (a *APIServiceImpl) StopWorkSession(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error) {
	return a.storage.StopWorkSession(employeeID)
}

func (a *APIServiceImpl) GetWorkSessions(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]entity.WorkSession, error) {
	return a.storage.GetWorkSessions(employeeID, date)
}
