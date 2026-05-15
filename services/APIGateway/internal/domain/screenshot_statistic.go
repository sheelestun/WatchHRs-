package domain

import (
	"time"

	"github.com/google/uuid"
)

// ScreenshotStatistic domain
type ScreenshotStatistic struct {
	ID                 uuid.UUID `json:"id" db:"id" validate:"required,uuid"`
	EmployeeID         uuid.UUID `json:"employeeID" db:"employee_id" validate:"required,uuid"`
	CountMouseClicks   int       `json:"count_mouse_clicks" db:"cnt_mouse_clicks" validate:"required,min=0"`
	CountKeyboardClick int       `json:"count_keyboard_clicks" db:"cnt_keyboard_clicks" validate:"required,min=0"`
	CreatedAt          time.Time `json:"timestamp" db:"created_at" validate:"required"`
}
