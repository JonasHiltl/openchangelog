package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/btvoidx/mint"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/events"
	"github.com/jonashiltl/openchangelog/internal/handler/rest"
	"github.com/jonashiltl/openchangelog/internal/handler/rss"
	"github.com/jonashiltl/openchangelog/internal/handler/web"
	"github.com/jonashiltl/openchangelog/internal/handler/web/admin"
	"github.com/jonashiltl/openchangelog/internal/load"
	"github.com/jonashiltl/openchangelog/internal/parse"
	"github.com/jonashiltl/openchangelog/internal/search"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/jonashiltl/openchangelog/internal/xcache"
	"github.com/jonashiltl/openchangelog/internal/xlog"
	"github.com/naveensrinivasan/httpcache"
	"github.com/rs/cors"
)

func main() {
	cfg, err := parseConfig()
	if err != nil {
		slog.Error("failed to read config", xlog.ErrAttr(err))
		os.Exit(1)
	}
	slog.SetDefault(xlog.NewLogger(cfg))

	mux := http.NewServeMux()
	cache, err := createCache(cfg)
	if err != nil {
		slog.Error("failed to create cache", xlog.ErrAttr(err))
		os.Exit(1)
	}

	st, err := createStore(cfg)
	if err != nil {
		slog.Error("failed to create store", xlog.ErrAttr(err))
		os.Exit(1)
	}

	searcher, err := createSearcher(cfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer searcher.Close()

	e := new(mint.Emitter)
	parser := parse.NewParser(parse.CreateGoldmark())
	loader := load.NewLoader(cfg, st, cache, parser, e)
	renderer := web.NewRenderer(cfg)
	listener := events.NewListener(cfg, e, parser, searcher, cache)
	listener.Start()
	defer listener.Close()

	rest.RegisterRestHandler(mux, rest.NewEnv(st, loader, parser, e))
	web.RegisterWebHandler(mux, web.NewEnv(cfg, loader, parser, renderer, searcher))
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
			return xcache.NewMemoryCache(), nil
		case config.Disk:
			if cfg.Cache.Disk == nil {
				return nil, errors.New("missing 'cache.file' config section")
			}
			slog.Info("using disk cache")
			return xcache.NewDiskCache(cfg), nil
		case config.S3:
			if cfg.Cache.S3 == nil {
				return nil, errors.New("missing 'cache.s3' config section")
			}
			slog.Info("using s3 cache")
			return xcache.NewS3Cache(cfg.Cache.S3.Bucket), nil
		}
	}
	return nil, nil
}

func createSearcher(cfg config.Config) (search.Searcher, error) {
	if cfg.Search == nil {
		slog.Debug("no search configuration defined, using noop searcher")
		return search.NewNoopSearcher(), nil
	}
	return search.NewSearcher(cfg)
}
