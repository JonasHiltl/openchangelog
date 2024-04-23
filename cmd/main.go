package main

import (
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/parse"
	"github.com/jonashiltl/openchangelog/internal/server"
	"github.com/jonashiltl/openchangelog/internal/source"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	p := parse.NewParser()
	s, err := source.NewFromConfig(cfg)
	if err != nil {
		panic(err)
	}

	srv := server.New(s, p, server.WithPort(80))
	srv.Start()
}
