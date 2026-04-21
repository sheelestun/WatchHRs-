package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/domain"
	"github.com/sheelestun/WatchHRs-/internal/web/handler"
)

type ManagerStorage interface {
	AddManager(ctx context.Context, manager domain.Manager) (uuid.UUID, error)
	RemoveManager(ctx context.Context, managerID uuid.UUID) error
}

type _ handler.ManagerService

type managerService struct {
	managerStorage ManagerStorage
	validate       *validator.Validate
}

func NewManagerService(managerStorage ManagerStorage, validate *validator.Validate) handler.ManagerService {
	return &managerService{
		managerStorage: managerStorage,
		validate:       validate,
	}
}

func (m *managerService) AddManager(ctx context.Context, manager domain.Manager) (uuid.UUID, error) {
	manager.ID = uuid.New()
	if err := m.validate.Struct(manager); err != nil {
		return uuid.Nil, err
	}
	return m.managerStorage.AddManager(ctx, manager)
}

func (m *managerService) RemoveManager(ctx context.Context, managerID uuid.UUID) error {
	return m.managerStorage.RemoveManager(ctx, managerID)
}
