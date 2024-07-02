package render

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/jonashiltl/openchangelog/components"
	"github.com/jonashiltl/openchangelog/internal/handler/web/views"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/jonashiltl/openchangelog/parse"
)

type Renderer interface {
	RenderIndex(ctx context.Context, w io.Writer, args RenderIndexArgs) error
}

type RenderIndexArgs struct {
	CL       store.Changelog
	Articles []parse.ParsedArticle
	HasMore  bool
	NextPage int
	PageSize int
}

func New() Renderer {
	return &renderer{}
}

type renderer struct{}

func (r *renderer) RenderIndex(ctx context.Context, w io.Writer, args RenderIndexArgs) error {
	articles := make([]components.ArticleArgs, len(args.Articles))
	for i, a := range args.Articles {
		buf := new(strings.Builder)
		_, err := io.Copy(buf, a.Content)
		if err != nil {
			continue
		}

		articles[i] = components.ArticleArgs{
			ID:          fmt.Sprint(a.Meta.PublishedAt.Unix()),
			Title:       a.Meta.Title,
			Description: a.Meta.Description,
			PublishedAt: a.Meta.PublishedAt.Format("02 Jan 2006"),
			Content:     buf.String(),
		}
	}

	return views.Index(views.IndexArgs{
		HeaderArgs: components.HeaderArgs{
			Title:    args.CL.Title,
			Subtitle: args.CL.Subtitle,
		},
		Logo: components.Logo{
			Src:    args.CL.LogoSrc,
			Width:  args.CL.LogoWidth,
			Height: args.CL.LogoHeight,
			Alt:    args.CL.LogoAlt,
			Link:   args.CL.LogoLink,
		},
		ArticleListArgs: components.ArticleListArgs{
			Articles: articles,
			HasMore:  args.HasMore,
			NextPage: args.NextPage,
			PageSize: args.PageSize,
		},
	}).Render(ctx, w)
}
