package web

import (
	"net/http"

	"github.com/jonashiltl/openchangelog/internal/adapters/web/routes"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/domain/changelog"
	"github.com/jonashiltl/openchangelog/parse"
)

type NewServerArgs struct {
	Cfg      config.Config
	Parser   parse.Parser
	CService changelog.Service
}

func NewServer(args NewServerArgs) http.Handler {
	mux := http.NewServeMux()

	e := routes.NewEnv(args.Cfg, args.Parser, args.CService)
	routes.Init(mux, e)

	return mux
}
