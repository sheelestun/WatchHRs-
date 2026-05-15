package domain

import "github.com/google/uuid"

// Manager domain
type Manager struct {
	ID    uuid.UUID `json:"id" db:"id" validate:"uuid"`
	Name  string    `json:"name" db:"name" validate:"required,min=2,max=32"`
	Email string    `json:"email" db:"email" validate:"required,email"`
}
