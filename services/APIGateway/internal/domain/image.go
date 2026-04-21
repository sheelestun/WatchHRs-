package domain

import "github.com/google/uuid"

// Image domain
type Image struct {
	ID     uuid.UUID `json:"id" bd:"id" validate:"uuid"`
	UserID uuid.UUID `json:"userID" bd:"user_id" validate:"uuid"`
}
