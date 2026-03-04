package service

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/entity"
)

// Storage interface for repository
type Storage interface {
	FindUser(ctx context.Context, userId uuid.UUID) (string, error)

	AddManager(ctx context.Context, manager entity.Manager) (uuid.UUID, error)
	RemoveManager(ctx context.Context, managerID uuid.UUID) error

	AddEmployee(ctx context.Context, employee entity.Employee) (uuid.UUID, error)
	GetAllEmployeesByManagerID(ctx context.Context, managerID uuid.UUID) ([]entity.Employee, error)
	RemoveEmployee(ctx context.Context, employeeID uuid.UUID) error

	AddPhoto(ctx context.Context, photo entity.Photo) (uuid.UUID, error)

	AddScreenshotStatistic(ctx context.Context, screenshot entity.ScreenshotStatistic) (uuid.UUID, error)
	GetScreenshotsStatistic(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]entity.ScreenshotStatistic, error)

	StartWorkSession(ctx context.Context, session entity.WorkSession) (uuid.UUID, error)
	StopWorkSession(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error)
	GetWorkSessions(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]entity.WorkSession, error)
}

type Cache interface {
	SaveTokenInCache(ctx context.Context, tokenID, userID string, expiresAt time.Time) error
	ExistsTokenInCache(ctx context.Context, tokenID string) (bool, error)
	DeleteTokenInCache(ctx context.Context, tokenID string) error
}

// APIServiceImpl implementation for ApiService interface
type APIServiceImpl struct {
	storage Storage
	cache   Cache

	validate *validator.Validate
}

// NewAPIService constructor for NewAPIService
func NewAPIService(storage Storage, cache Cache) *APIServiceImpl {
	return &APIServiceImpl{storage: storage,
		cache:    cache,
		validate: validator.New()}
}

func (a *APIServiceImpl) AddManager(ctx context.Context, manager entity.Manager) (uuid.UUID, error) {
	manager.ID = uuid.New()
	if err := a.validate.Struct(manager); err != nil {
		return uuid.Nil, err
	}
	return a.storage.AddManager(ctx, manager)
}

func (a *APIServiceImpl) RemoveManager(ctx context.Context, managerID uuid.UUID) error {
	return a.storage.RemoveManager(ctx, managerID)
}

func (a *APIServiceImpl) AddEmployee(ctx context.Context, employee entity.Employee) (uuid.UUID, error) {
	employee.ID = uuid.New()
	if err := a.validate.Struct(employee); err != nil {
		return uuid.Nil, err
	}
	return a.storage.AddEmployee(ctx, employee)
}

func (a *APIServiceImpl) GetEmployeesByManagerID(ctx context.Context, managerID uuid.UUID) ([]entity.Employee, error) {
	return a.storage.GetAllEmployeesByManagerID(ctx, managerID)
}

func (a *APIServiceImpl) RemoveEmployee(ctx context.Context, employeeID uuid.UUID) error {
	return a.storage.RemoveEmployee(ctx, employeeID)
}

func (a *APIServiceImpl) AddPhoto(ctx context.Context, photo entity.Photo) (uuid.UUID, error) {
	photo.ID = photo.UserID
	if err := a.validate.Struct(photo); err != nil {
		return uuid.Nil, err
	}
	return a.storage.AddPhoto(ctx, photo)
}

func (a *APIServiceImpl) AddScreenshotStatistic(ctx context.Context, screenshot entity.ScreenshotStatistic) (uuid.UUID, error) {
	screenshot.ID = uuid.New()
	screenshot.Timestamp = time.Now()
	if err := a.validate.Struct(screenshot); err != nil {
		return uuid.Nil, err
	}
	return a.storage.AddScreenshotStatistic(ctx, screenshot)
}

func (a *APIServiceImpl) GetScreenshotsStatistic(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]entity.ScreenshotStatistic, error) {
	return a.storage.GetScreenshotsStatistic(ctx, employeeID, date)
}

func (a *APIServiceImpl) StartWorkSession(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error) {
	newSession := entity.WorkSession{ID: uuid.New(), EmployeeID: employeeID, StartTime: time.Now()}
	if err := a.validate.Struct(newSession); err != nil {
		return uuid.Nil, err
	}
	return a.storage.StartWorkSession(ctx, newSession)
}

func (a *APIServiceImpl) StopWorkSession(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error) {
	return a.storage.StopWorkSession(ctx, employeeID)
}

func (a *APIServiceImpl) GetWorkSessions(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]entity.WorkSession, error) {
	return a.storage.GetWorkSessions(ctx, employeeID, date)
}

func (a *APIServiceImpl) Auth(ctx context.Context, userId uuid.UUID) (string, error) {
	return a.storage.FindUser(ctx, userId)
}

func (a *APIServiceImpl) DeleteToken(ctx context.Context, tokenID string) error {
	return a.cache.DeleteTokenInCache(ctx, tokenID)
}

func (a *APIServiceImpl) ExistsToken(ctx context.Context, tokenID string) (bool, error) {
	return a.cache.ExistsTokenInCache(ctx, tokenID)
}

func (a *APIServiceImpl) SaveToken(ctx context.Context, tokenID, userID string, expiresAt time.Time) error {
	return a.cache.SaveTokenInCache(ctx, tokenID, userID, expiresAt)
}
