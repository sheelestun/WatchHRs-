package domain

import "github.com/google/uuid"

// Image domain
type Image struct {
	ID     uuid.UUID `json:"id" db:"id" validate:"uuid"`
	UserID uuid.UUID `json:"userID" db:"user_id" validate:"uuid"`
}
