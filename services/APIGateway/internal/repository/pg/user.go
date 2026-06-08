package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (p *PostgresRepository) FindUser(ctx context.Context, userId uuid.UUID) (string, error) {
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
