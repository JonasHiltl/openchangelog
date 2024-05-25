package routes

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/jonashiltl/openchangelog/components"
	"github.com/jonashiltl/openchangelog/internal/adapters/web/views"
	"github.com/jonashiltl/openchangelog/loader"
)

func index(e *env, w http.ResponseWriter, r *http.Request) error {
	wId := r.URL.Query().Get("workspace-id")
	cIdStr := r.URL.Query().Get("changelog-id")
	cId, err := strconv.ParseInt(cIdStr, 10, 64)
	if cIdStr != "" && err != nil {
		return err
	}

	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page-size")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10
	}

	cl, err := e.changelogSrv.GetChangelog(r.Context(), wId, cId)
	if err != nil {
		return err
	}

	articles := make([]components.ArticleArgs, 0, pageSize)
	hasMore := false
	if cl.Source != nil {
		load, err := cl.Source.ToLoader(e.cfg)
		if err != nil {
			return err
		}

		res, err := e.parser.Parse(r.Context(), load, loader.NewPagination(pageSize, page))
		if err != nil {
			return err
		}

		hasMore = res.HasMore
		for _, a := range res.Articles {
			buf := new(strings.Builder)
			_, err := io.Copy(buf, a.Content)
			if err != nil {
				continue
			}

			articles = append(articles, components.ArticleArgs{
				ID:          fmt.Sprint(a.Meta.PublishedAt.Unix()),
				Title:       a.Meta.Title,
				Description: a.Meta.Description,
				PublishedAt: a.Meta.PublishedAt.Format("02 Jan 2006"),
				Content:     buf.String(),
			})
		}
	}

	articleListArgs := components.ArticleListArgs{
		Articles: articles,
		PageSize: pageSize,
		NextPage: page + 1,
		HasMore:  hasMore,
	}

	if htmxHeader := r.Header.Get("HX-Request"); len(htmxHeader) > 0 {
		if len(articles) > 0 {
			return components.ArticleList(articleListArgs).Render(r.Context(), w)
		} else {
			w.WriteHeader(http.StatusNoContent)
			return nil
		}
	}

	indexArgs := views.IndexArgs{
		HeaderArgs: components.HeaderArgs{
			Title:    cl.Title,
			Subtitle: cl.Subtitle,
		},
		NavbarArgs: components.NavbarArgs{
			Logo: components.Logo{
				Src:    cl.Logo.Src,
				Width:  cl.Logo.Width,
				Height: cl.Logo.Height,
				Alt:    cl.Logo.Alt,
				Link:   cl.Logo.Link,
			},
		},
		ArticleListArgs: articleListArgs,
	}

	return views.Index(indexArgs).Render(r.Context(), w)
}
