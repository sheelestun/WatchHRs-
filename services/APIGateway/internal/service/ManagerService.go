package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/domain"
)

type ManagerStorage interface {
	AddManager(ctx context.Context, manager domain.Manager) (uuid.UUID, error)
	RemoveManager(ctx context.Context, managerID uuid.UUID) error
}

type ManagerServiceImpl struct {
	managerStorage ManagerStorage
	validate       *validator.Validate
}

func NewManagerServiceImpl(managerStorage ManagerStorage, validate *validator.Validate) *ManagerServiceImpl {
	return &ManagerServiceImpl{
		managerStorage: managerStorage,
		validate:       validate,
	}
}

func (m *ManagerServiceImpl) AddManager(ctx context.Context, manager domain.Manager) (uuid.UUID, error) {
	manager.ID = uuid.New()
	if err := m.validate.Struct(manager); err != nil {
		return uuid.Nil, err
	}
	return m.managerStorage.AddManager(ctx, manager)
}

func (m *ManagerServiceImpl) RemoveManager(ctx context.Context, managerID uuid.UUID) error {
	return m.managerStorage.RemoveManager(ctx, managerID)
}
