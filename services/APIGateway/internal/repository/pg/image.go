package pg

import (
	"context"

	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/domain"
	log "github.com/sirupsen/logrus"
)

func (p *PostgresRepository) AddImage(ctx context.Context, photo domain.Image) (uuid.UUID, error) {
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
