package admin

import (
	"net/http"

	"github.com/jonashiltl/openchangelog/internal/handler"
	adminviews "github.com/jonashiltl/openchangelog/internal/handler/web/admin/views"
	"github.com/jonashiltl/openchangelog/internal/store"
	"golang.org/x/sync/errgroup"
)

func details(e *env, w http.ResponseWriter, r *http.Request) error {
	authorize := r.URL.Query().Get(handler.AUTHORIZE_QUERY)
	err := handler.ValidatePassword(e.cfg.Admin.PasswordHash, authorize)
	if err != nil {
		return err
	}

	wid, err := store.ParseWID(r.PathValue("wid"))
	if err != nil {
		return err
	}

	var ws store.Workspace
	var cls []store.Changelog
	var eg errgroup.Group

	eg.Go(func() error {
		ws, err = e.st.GetWorkspace(r.Context(), wid)
		return err
	})
	eg.Go(func() error {
		cls, err = e.st.ListChangelogs(r.Context(), wid)
		return err
	})

	err = eg.Wait()
	if err != nil {
		return err
	}

	return adminviews.WorkspaceDetails(adminviews.WorkspaceDetailsArgs{
		Workspace:  ws,
		Changelogs: cls,
	}).Render(r.Context(), w)
}
