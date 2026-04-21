package domain

import (
	"time"

	"github.com/google/uuid"
)

// ScreenshotStatistic domain
type ScreenshotStatistic struct {
	ID                 uuid.UUID `json:"id" bd:"id" validate:"required,uuid"`
	EmployeeID         uuid.UUID `json:"employeeID" bd:"employee_id" validate:"required,uuid"`
	CountMouseClicks   int       `json:"count_mouse_clicks" bd:"cnt_mouse_clicks" validate:"required,min=0"`
	CountKeyboardClick int       `json:"count_keyboard_clicks" bd:"cnt_keyboard_clicks" validate:"required,min=0"`
	CreatedAt          time.Time `json:"timestamp" bd:"created_at" validate:"required"`
}
