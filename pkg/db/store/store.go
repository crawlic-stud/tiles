package db

import (
	"context"
	"tiles/pkg/db"
	"tiles/pkg/db/gen"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool

	*gen.Queries
}

func New(ctx context.Context, dsn string) (*Store, error) {
	pool, err := db.NewPool(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &Store{
		pool:    pool,
		Queries: gen.New(pool),
	}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

// Transaction helper
func (s *Store) WithTx(ctx context.Context, fn func(*gen.Queries) error) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	// safe to call after commit
	defer tx.Rollback(ctx)

	if err := fn(s.Queries.WithTx(tx)); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
