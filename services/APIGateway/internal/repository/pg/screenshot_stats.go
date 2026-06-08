package pg

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/domain"
	log "github.com/sirupsen/logrus"
)

func (p *PostgresRepository) AddScreenshotStatistic(ctx context.Context, screenshot domain.ScreenshotStatistic) (uuid.UUID, error) {
	query := `
		INSERT INTO screenshot_statistics (id, employee_id, cnt_mouse_clicks, cnt_keyboard_clicks, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`
	var id uuid.UUID
	err := p.db.GetContext(ctx, &id, query, screenshot.ID, screenshot.EmployeeID, screenshot.CountMouseClicks, screenshot.CountKeyboardClick, screenshot.CreatedAt)
	if err != nil {
		return uuid.Nil, err
	}
	log.Debugf("Added screenshot %+v", screenshot)
	return id, nil
}

func (p *PostgresRepository) GetScreenshotsStatistic(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]domain.ScreenshotStatistic, error) {
	query := `
	SELECT id, employee_id, cnt_mouse_clicks, cnt_keyboard_clicks, created_at
	FROM screenshot_statistics
	WHERE employee_id = $1
	  AND created_at >= $2
	  AND created_at <  $2 + INTERVAL '1 day'
`

	startOfDay := time.Date(
		date.Year(), date.Month(), date.Day(),
		0, 0, 0, 0, date.Location(),
	)

	var screenshots []domain.ScreenshotStatistic
	err := p.db.SelectContext(ctx, &screenshots, query, employeeID, startOfDay)
	if err != nil {
		return nil, err
	}
	log.Debugf("Get screenshots %+v", screenshots)
	return screenshots, nil
}
