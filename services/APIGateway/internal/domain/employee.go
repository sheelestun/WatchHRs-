package domain

import "github.com/google/uuid"

// Employee domain
type Employee struct {
	ID        uuid.UUID `json:"id" bd:"id" validate:"uuid"`
	Name      string    `json:"name" bd:"name" validate:"required,min=2,max=32"`
	Email     string    `json:"email" bd:"email" validate:"required,email"`
	ManagerID uuid.UUID `json:"managerID" bd:"manager_id" validate:"required,uuid"`
}
