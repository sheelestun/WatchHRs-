package domain

import "github.com/google/uuid"

// Manager domain
type Manager struct {
	ID    uuid.UUID `json:"id" bd:"id" validate:"uuid"`
	Name  string    `json:"name" bd:"name" validate:"required,min=2,max=32"`
	Email string    `json:"email" bd:"email" validate:"required,email"`
}
