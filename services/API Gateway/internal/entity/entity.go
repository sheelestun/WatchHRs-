package entity

import (
	"time"

	"github.com/google/uuid"
)

type Manager struct {
	Id    uuid.UUID `json:"id" bd:"id" validate:"uuid"`
	Name  string    `json:"name" bd:"name" validate:"required,min=2,max=32"`
	Email string    `json:"email" bd:"email" validate:"required,email"`
}

type Employee struct {
	Id        uuid.UUID `json:"id" bd:"id" validate:"uuid"`
	Name      string    `json:"name" bd:"name" validate:"required,min=2,max=32"`
	Email     string    `json:"email" bd:"email" validate:"required,email"`
	ManagerId uuid.UUID `json:"manager_id" bd:"manager_id" validate:"required,uuid"`
}

type Photo struct {
	Id     uuid.UUID `json:"id" bd:"id" validate:"uuid"`
	UserId uuid.UUID `json:"userId" bd:"userId" validate:"uuid"`
}

type ScreenshotStatistic struct {
	Id                 uuid.UUID `json:"id" bd:"id" validate:"uuid"`
	EmployeeId         uuid.UUID `json:"employeeId" bd:"employeeId" validate:"required,uuid"`
	CountMouseClicks   int       `json:"count_mouse_clicks" bd:"cnt_mouse_clicks" validate:"required,min=0"`
	CountKeyboardClick int       `json:"count_keyboard_clicks" bd:"cnt_keyboard_clicks" validate:"required,min=0"`
	Timestamp          time.Time `json:"timestamp" bd:"timestamp" validate:"required"`
}

type WorkSession struct {
	Id         uuid.UUID     `json:"id" bd:"id" validate:"uuid"`
	EmployeeId uuid.UUID     `json:"employee" bd:"employee" validate:"uuid"`
	StartTime  time.Time     `json:"start_time" bd:"start_time" validate:"required"`
	EndTime    time.Time     `json:"end_time" bd:"end_time"`
	TotalTime  time.Duration `json:"total_time" bd:"total_time"`
}
