package rest

import (
	"net/http"

	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/domain/changelog"
	"github.com/jonashiltl/openchangelog/internal/domain/source"
	"github.com/jonashiltl/openchangelog/internal/domain/workspace"
)

type NewServerArgs struct {
	Cfg          config.Config
	SourceSrv    source.Service
	WorkspaceSrv workspace.Service
	ChangelogSrv changelog.Service
}

func NewServer(args NewServerArgs) http.Handler {
	mux := http.NewServeMux()

	e := env{
		cfg:          args.Cfg,
		sourceSrv:    args.SourceSrv,
		workspaceSrv: args.WorkspaceSrv,
		changelogSrv: args.ChangelogSrv,
	}
	initRoutes(mux, &e)

	return mux
}
