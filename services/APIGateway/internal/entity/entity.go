package entity

import (
	"time"

	"github.com/google/uuid"
)

// Manager entity
type Manager struct {
	ID    uuid.UUID `json:"id" bd:"id" validate:"uuid"`
	Name  string    `json:"name" bd:"name" validate:"required,min=2,max=32"`
	Email string    `json:"email" bd:"email" validate:"required,email"`
}

// Employee entity
type Employee struct {
	ID        uuid.UUID `json:"id" bd:"id" validate:"uuid"`
	Name      string    `json:"name" bd:"name" validate:"required,min=2,max=32"`
	Email     string    `json:"email" bd:"email" validate:"required,email"`
	ManagerID uuid.UUID `json:"managerID" bd:"manager_id" validate:"required,uuid"`
}

// Photo entity
type Photo struct {
	ID     uuid.UUID `json:"id" bd:"id" validate:"uuid"`
	UserID uuid.UUID `json:"userID" bd:"user_id" validate:"uuid"`
}

// ScreenshotStatistic entity
type ScreenshotStatistic struct {
	ID                 uuid.UUID `json:"id" bd:"id" validate:"uuid"`
	EmployeeID         uuid.UUID `json:"employeeID" bd:"employee_id" validate:"required,uuid"`
	CountMouseClicks   int       `json:"count_mouse_clicks" bd:"cnt_mouse_clicks" validate:"required,min=0"`
	CountKeyboardClick int       `json:"count_keyboard_clicks" bd:"cnt_keyboard_clicks" validate:"required,min=0"`
	Timestamp          time.Time `json:"timestamp" bd:"created_at" validate:"required"`
}

// WorkSession entity
type WorkSession struct {
	ID         uuid.UUID     `json:"id" bd:"id" validate:"uuid"`
	EmployeeID uuid.UUID     `json:"employeeID" bd:"employee_id" validate:"uuid"`
	StartTime  time.Time     `json:"start_time" bd:"start_time" validate:"required"`
	EndTime    time.Time     `json:"end_time" bd:"end_time"`
	TotalTime  time.Duration `json:"total_time" bd:"total_time"`
}
