package web

import (
	"errors"
	"net/http"

	"github.com/jonashiltl/openchangelog/components"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/search"
	"github.com/jonashiltl/openchangelog/internal/source"
)

func searchSubmit(e *env, w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return errs.NewBadRequest(err)
	}

	cl, err := e.loader.GetChangelog(r)
	if err != nil {
		return errs.NewBadRequest(err)
	}

	if cl.Protected {
		err = ensurePasswordProvided(r, cl.PasswordHash)
		if err != nil {
			return errs.NewUnauthorized(err)
		}
	}

	sid := source.NewIDFromChangelog(cl)
	if sid == "" {
		return errs.NewBadRequest(errors.New("changelog has no active source"))
	}

	q := r.FormValue("query")
	if q == "" {
		return components.SearchResults(components.SearchResultsArgs{
			Result: search.SearchResults{},
		}).Render(r.Context(), w)
	}

	res, err := e.searcher.Search(r.Context(), search.SearchArgs{
		SID:   sid.String(),
		Query: q,
	})
	if err != nil {
		return errs.NewBadRequest(err)
	}

	return components.SearchResults(components.SearchResultsArgs{
		Query:  q,
		Result: res,
	}).Render(r.Context(), w)
}
