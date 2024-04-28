package api

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jonashiltl/openchangelog/internal/store"
)

type api struct {
	queries *store.Queries
	pool    *pgxpool.Pool
}

func New(queries *store.Queries, pool *pgxpool.Pool) *api {
	return &api{
		queries: queries,
		pool:    pool,
	}
}
