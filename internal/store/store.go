package store

import (
	"context"

	db "github.com/KrishnaGrg1/pulseway/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{
		Pool:    pool,
		Queries: db.New(pool),
	}
}

func Connect(dbUrl string) (*Store, error) {

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return New(pool), nil
}
