package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/domain"
	"github.com/sheelestun/WatchHRs-/internal/web/handler"
)

type EmployeeStorage interface {
	AddEmployee(ctx context.Context, employee domain.Employee) (uuid.UUID, error)
	GetAllEmployeesByManagerID(ctx context.Context, managerID uuid.UUID) ([]domain.Employee, error)
	RemoveEmployee(ctx context.Context, employeeID uuid.UUID) error
}

type _ handler.EmployeeService

type employeeService struct {
	employeeStorage EmployeeStorage
	validate        *validator.Validate
}

func NewEmployeeService(employeeStorage EmployeeStorage, validate *validator.Validate) handler.EmployeeService {
	return &employeeService{
		employeeStorage: employeeStorage,
		validate:        validate,
	}
}

func (e *employeeService) AddEmployee(ctx context.Context, employee domain.Employee) (uuid.UUID, error) {
	employee.ID = uuid.New()
	if err := e.validate.Struct(employee); err != nil {
		return uuid.Nil, err
	}
	return e.employeeStorage.AddEmployee(ctx, employee)
}

func (e *employeeService) GetEmployeesByManagerID(ctx context.Context, managerID uuid.UUID) ([]domain.Employee, error) {
	return e.employeeStorage.GetAllEmployeesByManagerID(ctx, managerID)
}

func (e *employeeService) RemoveEmployee(ctx context.Context, employeeID uuid.UUID) error {
	return e.employeeStorage.RemoveEmployee(ctx, employeeID)
}
