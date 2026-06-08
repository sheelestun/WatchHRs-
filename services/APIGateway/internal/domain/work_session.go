package domain

import (
	"time"

	"github.com/google/uuid"
)

// WorkSession domain
type WorkSession struct {
	ID         uuid.UUID  `json:"id" db:"id" validate:"uuid"`
	EmployeeID uuid.UUID  `json:"employeeID" db:"employee_id" validate:"uuid"`
	StartTime  time.Time  `json:"start_time" db:"start_time" validate:"required"`
	EndTime    *time.Time `json:"end_time,omitempty" db:"end_time"`
	TotalTime  *string    `json:"total_time,omitempty" db:"total_time"`
}
