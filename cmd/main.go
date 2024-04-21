package main

import (
	"github.com/jonashiltl/openchangelog/internal/parse"
	"github.com/jonashiltl/openchangelog/internal/server"
	"github.com/jonashiltl/openchangelog/internal/source"
)

func main() {
	p := parse.NewParser()
	s := source.LocalFileSource(".testdata")
	srv := server.New(s, p, server.WithPort(4000))
	srv.Start()
}
