package pg

import (
	"github.com/jmoiron/sqlx"
)

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPgRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}
