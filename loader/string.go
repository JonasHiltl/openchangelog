package loader

import (
	"context"
	"io"
	"strings"
)

type stringSource struct {
	articles []string
}

func String(articles []string) Loader {
	return stringSource{
		articles: articles,
	}
}

func (s stringSource) Load(ctx context.Context, page Pagination) (LoadResult, error) {
	if page.StartIdx() >= len(s.articles) {
		return LoadResult{
			Articles: []RawArticle{},
			HasMore:  false,
		}, nil
	}

	articles := make([]RawArticle, 0, page.PageSize())
	for i := page.StartIdx(); i <= page.EndIdx() && i < len(s.articles); i++ {
		articles = append(articles, RawArticle{
			Content: io.NopCloser(strings.NewReader(s.articles[i])),
		})
	}

	return LoadResult{
		Articles: articles,
		HasMore:  page.EndIdx()+1 < len(articles),
	}, nil
}
