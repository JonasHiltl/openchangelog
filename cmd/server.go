package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/handler/rest"
	"github.com/jonashiltl/openchangelog/internal/handler/rss"
	"github.com/jonashiltl/openchangelog/internal/handler/web"
	"github.com/jonashiltl/openchangelog/internal/handler/web/admin"
	"github.com/jonashiltl/openchangelog/internal/lgr"
	"github.com/jonashiltl/openchangelog/internal/load"
	"github.com/jonashiltl/openchangelog/internal/parse"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/naveensrinivasan/httpcache"
	"github.com/naveensrinivasan/httpcache/diskcache"
	"github.com/peterbourgon/diskv"
	"github.com/rs/cors"
	"github.com/sourcegraph/s3cache"
)

func main() {
	cfg, err := parseConfig()
	if err != nil {
		slog.Error("failed to read config", lgr.ErrAttr(err))
		os.Exit(1)
	}
	slog.SetDefault(lgr.NewLogger(cfg))

	mux := http.NewServeMux()
	cache, err := createCache(cfg)
	if err != nil {
		slog.Error("failed to create cache", lgr.ErrAttr(err))
		os.Exit(1)
	}

	st, err := createStore(cfg)
	if err != nil {
		slog.Error("failed to create store", lgr.ErrAttr(err))
		os.Exit(1)
	}

	loader := load.NewLoader(cfg, st, cache)
	parser := parse.NewParser(parse.CreateGoldmark())
	renderer := web.NewRenderer(cfg)

	rest.RegisterRestHandler(mux, rest.NewEnv(st, loader, parser))
	web.RegisterWebHandler(mux, web.NewEnv(cfg, loader, parser, renderer))
	admin.RegisterAdminHandler(mux, admin.NewEnv(cfg, st))
	rss.RegisterRSSHandler(mux, rss.NewEnv(cfg, loader, parser))
	handler := cors.Default().Handler(mux)

	slog.Info("Ready to serve requests", slog.String("addr", fmt.Sprintf("http://%s", cfg.Addr)))
	log.Fatal(http.ListenAndServe(cfg.Addr, handler))
}

func parseConfig() (config.Config, error) {
	configPath := flag.String("config", "", "config file path")
	flag.Parse()
	return config.Load(*configPath)
}

func createStore(cfg config.Config) (store.Store, error) {
	if cfg.IsDBMode() {
		slog.Info("Starting Openchangelog backed by sqlite")
		return store.NewSQLiteStore(cfg.SqliteURL)
	} else {
		slog.Info("Starting Openchangelog in config mode")
		return store.NewConfigStore(cfg), nil
	}
}

func createCache(cfg config.Config) (httpcache.Cache, error) {
	if cfg.Cache != nil {
		switch cfg.Cache.Type {
		case config.Memory:
			slog.Info("using memory cache")
			return httpcache.NewMemoryCache(), nil
		case config.Disk:
			if cfg.Cache.Disk == nil {
				return nil, errors.New("missing 'cache.file' config section")
			}
			slog.Info("using disk cache")
			return diskcache.NewWithDiskv(diskv.New(diskv.Options{
				BasePath:     cfg.Cache.Disk.Location,
				CacheSizeMax: cfg.Cache.Disk.MaxSize, // bytes
			})), nil
		case config.S3:
			if cfg.Cache.S3 == nil {
				return nil, errors.New("missing 'cache.s3' config section")
			}
			slog.Info("using s3 cache")
			return s3cache.New(cfg.Cache.S3.Bucket), nil
		}
	}
	return nil, nil
}
