package admin

import (
	"net/http"

	"github.com/jonashiltl/openchangelog/internal/handler"
	adminviews "github.com/jonashiltl/openchangelog/internal/handler/web/admin/views"
	"github.com/jonashiltl/openchangelog/internal/handler/web/static"
)

func adminOverview(e *env, w http.ResponseWriter, r *http.Request) error {
	authorize := r.URL.Query().Get(handler.AUTHORIZE_QUERY)
	err := handler.ValidatePassword(e.cfg.Admin.PasswordHash, authorize)
	if err != nil {
		return err
	}

	rows, err := e.st.ListWorkspacesChangelogCount(r.Context())
	if err != nil {
		return err
	}

	return adminviews.Overview(adminviews.OverviewArgs{
		CSS:        static.AdminCSS,
		Workspaces: rows,
		Authorize:  authorize,
	}).Render(r.Context(), w)
}
