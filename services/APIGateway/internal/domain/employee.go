package domain

import "github.com/google/uuid"

// Employee domain
type Employee struct {
	ID        uuid.UUID `json:"id" db:"id" validate:"uuid"`
	Name      string    `json:"name" db:"name" validate:"required,min=2,max=32"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	ManagerID uuid.UUID `json:"managerID" db:"manager_id" validate:"required,uuid"`
}
