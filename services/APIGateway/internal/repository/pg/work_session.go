package pg

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/domain"
	log "github.com/sirupsen/logrus"
)

func (p *PostgresRepository) StartWorkSession(ctx context.Context, session domain.WorkSession) (uuid.UUID, error) {
	query := `
		INSERT INTO work_sessions (id, employee_id, start_time)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	var id uuid.UUID
	err := p.db.GetContext(ctx, &id, query,
		session.ID,
		session.EmployeeID,
		session.StartTime,
	)
	if err != nil {
		return uuid.Nil, errors.New("last work session did not stop yet")
	}

	return id, nil
}

func (p *PostgresRepository) StopWorkSession(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error) {
	query := `
		UPDATE work_sessions
		SET end_time = NOW(),
		    total_time = NOW() - start_time
		WHERE employee_id = $1
		  AND end_time IS NULL
		RETURNING id;
	`

	var id uuid.UUID
	err := p.db.GetContext(ctx, &id, query, employeeID)
	if err != nil {
		return uuid.Nil, errors.New("work session does not exist or already stopped")
	}

	return id, nil
}

func (p *PostgresRepository) GetWorkSessions(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]domain.WorkSession, error) {
	query := `
		SELECT id, employeeI_id, start_time, end_time, total_time
		FROM work_sessions
		WHERE employee_id = $1
		AND start_time >= $2
		AND start_time <  $2 + INTERVAL '1 day'
	`

	startOfDay := time.Date(
		date.Year(), date.Month(), date.Day(),
		0, 0, 0, 0, date.Location(),
	)

	var sessions []domain.WorkSession
	err := p.db.SelectContext(ctx, &sessions, query, employeeID, startOfDay)
	if err != nil {
		return nil, err
	}
	log.Debugf("Get work sessions %+v", sessions)
	return sessions, nil
}
