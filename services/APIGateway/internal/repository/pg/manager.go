package pg

import (
	"context"

	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/domain"
	log "github.com/sirupsen/logrus"
)

func (p *PostgresRepository) AddManager(ctx context.Context, manager domain.Manager) (uuid.UUID, error) {
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

func (p *PostgresRepository) RemoveManager(ctx context.Context, managerID uuid.UUID) error {
	query := `
		DELETE FROM managers WHERE id = $1
	`
	_, err := p.db.ExecContext(ctx, query, managerID)
	log.Debugf("Removed manager %s, with err %s", managerID, err.Error())
	return err
}
