package main

import (
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/parse"
	"github.com/jonashiltl/openchangelog/internal/server"
	"github.com/jonashiltl/openchangelog/internal/source"
)

func main() {
	_, err := config.Load()
	if err != nil {
		panic(err)
	}

	p := parse.NewParser()
	/* s, err := source.GitHub(source.GitHubSourceOptions{
		Owner:               "jonashiltl",
		Repository:          "openchangelog",
		Path:                ".testdata",
		GHAppPrivateKey:     cfg.GH_APP_PRIVATE_KEY,
		GHAppInstallationId: cfg.GH_APP_INSTALLATION_ID,
	})
	if err != nil {
		panic(err)
	}
	*/
	s := source.LocalFile(".testdata")

	srv := server.New(s, p, server.WithPort(80))
	srv.Start()
}
