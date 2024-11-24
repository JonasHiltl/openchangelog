package web

import (
	"errors"
	"net/http"
	"strings"

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

	if !cl.Searchable {
		return errs.NewBadRequest(errors.New("changelog is not searchable"))
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
	var tags []string
	for name, value := range r.PostForm {
		if strings.HasPrefix(name, "tag-") && len(value) > 0 && value[0] == "on" {
			tags = append(tags, strings.TrimPrefix(name, "tag-"))
		}
	}

	if q == "" && len(tags) == 0 {
		return components.SearchResults(components.SearchResultsArgs{
			Result: search.SearchResults{},
		}).Render(r.Context(), w)
	}

	res, err := e.searcher.Search(r.Context(), search.SearchArgs{
		SID:   sid.String(),
		Query: q,
		Tags:  tags,
	})
	if err != nil {
		return errs.NewBadRequest(err)
	}

	return components.SearchResults(components.SearchResultsArgs{
		Query:  q,
		Result: res,
	}).Render(r.Context(), w)
}

func searchTags(e *env, w http.ResponseWriter, r *http.Request) error {
	cl, err := e.loader.GetChangelog(r)
	if err != nil {
		return errs.NewBadRequest(err)
	}

	if !cl.Searchable {
		return errs.NewBadRequest(errors.New("changelog is not searchable"))
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

	tags := e.searcher.GetAllTags(r.Context(), sid.String())
	return components.TagSelectors(tags).Render(r.Context(), w)
}
