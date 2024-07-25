package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/handler/rest"
	"github.com/jonashiltl/openchangelog/internal/handler/web"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/jonashiltl/openchangelog/parse"
	"github.com/jonashiltl/openchangelog/render"
	"github.com/naveensrinivasan/httpcache"
	"github.com/naveensrinivasan/httpcache/diskcache"
	"github.com/peterbourgon/diskv"
	"github.com/sourcegraph/s3cache"
)

func main() {
	cfg, err := parseConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read config: %s\n", err)
		os.Exit(1)
	}
	mux := http.NewServeMux()
	cache, err := createCache(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create cache: %s\n", err)
		os.Exit(1)
	}

	st, err := createStore(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create store: %s\n", err)
		os.Exit(1)
	}

	rest.RegisterRestHandler(mux, rest.NewEnv(st))
	web.RegisterWebHandler(mux, web.NewEnv(cfg, st, render.New(), parse.NewParser(), cache))

	addr := fmt.Sprintf("localhost:%d", cfg.Port)
	fmt.Printf("Starting server at http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func parseConfig() (config.Config, error) {
	configPath := flag.String("config", "", "config file path")
	flag.Parse()
	return config.Load(*configPath)
}

func createStore(cfg config.Config) (store.Store, error) {
	if cfg.IsDBMode() {
		log.Println("Starting Openchangelog backed by sqlite")
		return store.NewSQLiteStore(cfg.SqliteURL)
	} else {
		log.Println("Starting Openchangelog in config mode")
		return store.NewConfigStore(cfg), nil
	}
}

func createCache(cfg config.Config) (httpcache.Cache, error) {
	if cfg.Cache != nil {
		switch cfg.Cache.Type {
		case config.Memory:
			log.Println("using memory cache")
			return httpcache.NewMemoryCache(), nil
		case config.Disk:
			if cfg.Cache.Disk == nil {
				return nil, errors.New("missing 'cache.file' config section")
			}
			log.Println("using disk cache")
			return diskcache.NewWithDiskv(diskv.New(diskv.Options{
				BasePath:     cfg.Cache.Disk.Location,
				CacheSizeMax: cfg.Cache.Disk.MaxSize, // bytes
			})), nil
		case config.S3:
			if cfg.Cache.S3 == nil {
				return nil, errors.New("missing 'cache.s3' config section")
			}
			log.Println("using s3 cache")
			return s3cache.New(cfg.Cache.S3.Bucket), nil
		}
	}
	return nil, nil
}
