package source

import (
	"context"
)

type stringSource struct {
	articles []string
}

func String(articles []string) Source {
	return stringSource{
		articles: articles,
	}
}

func (s stringSource) Load(ctx context.Context, params LoadParams) (LoadResult, error) {
	if params.StartIdx() >= len(s.articles) {
		return LoadResult{
			Articles: []Article{},
			HasMore:  false,
		}, nil
	}

	articles := make([]Article, 0, params.PageSize())
	for i := params.StartIdx(); i <= params.EndIdx() && i < len(s.articles); i++ {
		articles = append(articles, Article{
			Bytes: []byte(s.articles[i]),
		})
	}

	return LoadResult{
		Articles: articles,
		HasMore:  params.EndIdx()+1 < len(articles),
	}, nil
}
