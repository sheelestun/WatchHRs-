package domain

import (
	"time"

	"github.com/google/uuid"
)

// WorkSession domain
type WorkSession struct {
	ID         uuid.UUID     `json:"id" bd:"id" validate:"uuid"`
	EmployeeID uuid.UUID     `json:"employeeID" bd:"employee_id" validate:"uuid"`
	StartTime  time.Time     `json:"start_time" bd:"start_time" validate:"required"`
	EndTime    time.Time     `json:"end_time" bd:"end_time"`
	TotalTime  time.Duration `json:"total_time" bd:"total_time"`
}
