package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sheelestun/WatchHRs-/internal/entity"
	log "github.com/sirupsen/logrus"
)

type PgRepository struct {
	db *sqlx.DB
}

func NewPgRepository(db *sqlx.DB) *PgRepository {
	return &PgRepository{db: db}
}

func (p *PgRepository) FindUser(ctx context.Context, userId uuid.UUID) (string, error) {
	var role string
	query := `
		SELECT role FROM (
			SELECT 'manager' AS role FROM managers WHERE id = $1
			UNION ALL
			SELECT 'employee' AS role FROM employees WHERE id = $1
		) t
		LIMIT 1;
	`

	err := p.db.GetContext(ctx, &role, query, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	log.Debugf("Found user: %s", userId)
	return role, nil
}

func (p *PgRepository) AddManager(ctx context.Context, manager entity.Manager) (uuid.UUID, error) {
	query := `
		INSERT INTO managers (id, name, email)
		VALUES ($1, $2, $3)
		RETURNING id;
	`
	var id uuid.UUID
	err := p.db.GetContext(ctx, &id, query, manager.ID, manager.Name, manager.Email)
	if err != nil {
		return uuid.Nil, err
	}
	log.Debugf("Added manager %+v", id)
	return id, nil
}

func (p *PgRepository) RemoveManager(ctx context.Context, managerID uuid.UUID) error {
	query := `
		DELETE FROM managers WHERE id = $1
	`
	_, err := p.db.ExecContext(ctx, query, managerID)
	log.Debugf("Removed manager %s, with err %s", managerID, err.Error())
	return err
}

func (p *PgRepository) AddEmployee(ctx context.Context, employee entity.Employee) (uuid.UUID, error) {
	query := `
		INSERT INTO employees (id, name, email, manager_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	var id uuid.UUID
	err := p.db.GetContext(ctx, &id, query, employee.ID, employee.Name, employee.Email, employee.ManagerID)
	if err != nil {
		return uuid.Nil, err
	}
	log.Debugf("added employee in database: %+v", employee)
	return id, nil
}

func (p *PgRepository) GetAllEmployeesByManagerID(ctx context.Context, managerID uuid.UUID) ([]entity.Employee, error) {
	var employees []entity.Employee

	query := `
		SELECT id, name, email, manager_id
		FROM employees
		WHERE managerID = $1;
	`

	err := p.db.SelectContext(ctx, &employees, query, managerID)
	if err != nil {
		return nil, err
	}

	log.Debugf("Get all employees in database: %+v", employees)
	return employees, nil
}

func (p *PgRepository) RemoveEmployee(ctx context.Context, employeeID uuid.UUID) error {
	query := `
		DELETE FROM employees WHERE id = $1;
	`
	_, err := p.db.ExecContext(ctx, query, employeeID)
	log.Debugf("Removed employee %s, with err %s", employeeID, err.Error())
	return err

}

func (p *PgRepository) AddPhoto(ctx context.Context, photo entity.Photo) (uuid.UUID, error) {
	query := `
		INSERT INTO photos (id, user_id)
		VALUES ($1, $2)
	`
	var id uuid.UUID
	err := p.db.GetContext(ctx, &id, query, photo.ID, photo.UserID)
	if err != nil {
		return uuid.Nil, err
	}
	log.Debugf("Added photo %+v", photo)
	return id, nil
}

func (p *PgRepository) AddScreenshotStatistic(ctx context.Context, screenshot entity.ScreenshotStatistic) (uuid.UUID, error) {
	query := `
		INSERT INTO screenshot_statistics (id, employee_id, cnt_mouse_clicks, cnt_keyboard_clicks, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`
	var id uuid.UUID
	err := p.db.GetContext(ctx, &id, query, screenshot.ID, screenshot.EmployeeID, screenshot.CountMouseClicks, screenshot.CountKeyboardClick, screenshot.Timestamp)
	if err != nil {
		return uuid.Nil, err
	}
	log.Debugf("Added screenshot %+v", screenshot)
	return id, nil
}

func (p *PgRepository) GetScreenshotsStatistic(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]entity.ScreenshotStatistic, error) {
	query := `
	SELECT id, employee_id, cnt_mouse_clicks, cnt_keyboard_clicks, created_at
	FROM screenshot_statistics
	WHERE employee_id = $1
	WHERE employee_id = $1
	  AND created_at >= $2
	  AND created_at <  $2 + INTERVAL '1 day'
`

	startOfDay := time.Date(
		date.Year(), date.Month(), date.Day(),
		0, 0, 0, 0, date.Location(),
	)

	var screenshots []entity.ScreenshotStatistic
	err := p.db.SelectContext(ctx, &screenshots, query, employeeID, startOfDay)
	if err != nil {
		return nil, err
	}
	log.Debugf("Get screenshots %+v", screenshots)
	return screenshots, nil
}

func (p *PgRepository) StartWorkSession(ctx context.Context, session entity.WorkSession) (uuid.UUID, error) {
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

func (p *PgRepository) StopWorkSession(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error) {
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

func (p *PgRepository) GetWorkSessions(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]entity.WorkSession, error) {
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

	var sessions []entity.WorkSession
	err := p.db.SelectContext(ctx, &sessions, query, employeeID, startOfDay)
	if err != nil {
		return nil, err
	}
	log.Debugf("Get work sessions %+v", sessions)
	return sessions, nil
}
