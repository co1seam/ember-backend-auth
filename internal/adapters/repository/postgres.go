package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

type Postgres struct {
	db     *sql.DB
	cancel context.CancelFunc
}

func Init(ctx context.Context) (*Postgres, error) {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/postgres?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error pinging postgres: %w", err)
	}

	_, cancel := context.WithCancel(ctx)

	return &Postgres{db: db, cancel: cancel}, nil
}

func (pg *Postgres) Close() error {
	pg.cancel()
	return pg.db.Close()
}
