package main

import (
	"context"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/server"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/jonashiltl/openchangelog/parse"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read config: %v\n", err)
		os.Exit(1)
	}

	p := parse.NewParser()

	var queries *store.Queries
	var pool *pgxpool.Pool
	if cfg.DatabaseURL != "" {
		newPool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer newPool.Close()

		pool = newPool
		queries = store.New(newPool)
		m, err := migrate.New(
			"file://internal/store/migrations",
			cfg.DatabaseURL,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create migration instance: %v\n", err)
			os.Exit(1)
		}

		err = m.Up()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to run migrations: %v\n", err)
		}
	}

	srv := server.New(server.ServerArgs{
		Parser:  p,
		Cfg:     cfg,
		Queries: queries,
		Pool:    pool,
	})
	srv.Start()
}
