package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jonashiltl/openchangelog/internal/adapters/rest"
	"github.com/jonashiltl/openchangelog/internal/adapters/web"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/domain/changelog"
	"github.com/jonashiltl/openchangelog/internal/domain/source"
	"github.com/jonashiltl/openchangelog/internal/domain/workspace"
	"github.com/jonashiltl/openchangelog/parse"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read config: %v\n", err)
		os.Exit(1)
	}

	p := parse.NewParser()

	var pool *pgxpool.Pool
	if cfg.DatabaseURL != "" {
		newPool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer newPool.Close()

		pool = newPool
		m, err := migrate.New(
			"file://internal/migrations",
			cfg.DatabaseURL,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create migration instance: %v\n", err)
			os.Exit(1)
		}

		err = m.Up()
		if err != nil {
			if err != migrate.ErrNoChange {
				fmt.Fprintf(os.Stderr, "Failed to run migrations: %v\n", err)
			}
		}
	}

	var wRepo workspace.Repo
	var sRepo source.Repo
	var cRepo changelog.Repo
	if pool != nil {
		wRepo = workspace.NewPGRepo(pool)
		sRepo = source.NewPGRepo(pool)
		cRepo = changelog.NewPGRepo(pool)
	} else {
		wRepo = workspace.NewConfigRepo(cfg)
		sRepo = source.NewConfigRepo(cfg)
		cRepo = changelog.NewConfigRepo(cfg)
	}

	wService := workspace.NewService(wRepo)
	sService := source.NewService(sRepo)
	cService := changelog.NewService(cRepo)

	restHandler := rest.NewServer(rest.NewServerArgs{
		Cfg:          cfg,
		SourceSrv:    sService,
		WorkspaceSrv: wService,
		ChangelogSrv: cService,
	})

	webHandler := web.NewServer(web.NewServerArgs{
		Cfg:      cfg,
		Parser:   p,
		CService: cService,
	})

	mux := http.NewServeMux()
	mux.Handle("/", webHandler)
	mux.Handle("/api/{pathname...}", http.StripPrefix("/api", restHandler))

	addr := fmt.Sprintf("localhost:%d", cfg.Port)
	fmt.Printf("Starting server at http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
