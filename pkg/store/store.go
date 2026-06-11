package store

import (
	"context"
	"database/sql"
	"tiles/pkg/db"
	"tiles/pkg/server/cache"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct {
	DB        *sql.DB
	GameCache *cache.Cache[int64, db.Game]
	*db.Queries
}

func New(sqlitePath string) (*Store, error) {
	conn, err := sql.Open("sqlite", sqlitePath)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &Store{
		DB:        conn,
		Queries:   db.New(conn),
		GameCache: cache.New[int64, db.Game](time.Minute),
	}, nil
}

func (s *Store) Close() error {
	return s.DB.Close()
}

// Optional transaction helper
func (s *Store) WithTx(
	ctx context.Context,
	fn func(*db.Queries) error,
) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := s.Queries.WithTx(tx)

	if err := fn(q); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
