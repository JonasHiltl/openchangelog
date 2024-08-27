package render

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/jonashiltl/openchangelog/components"
	"github.com/jonashiltl/openchangelog/internal/changelog"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/handler/web/views"
	"github.com/jonashiltl/openchangelog/internal/handler/web/views/layout"
	"github.com/jonashiltl/openchangelog/internal/store"
)

type Renderer interface {
	RenderIndex(ctx context.Context, w io.Writer, args RenderIndexArgs) error
	RenderArticleList(ctx context.Context, w io.Writer, args RenderArticleListArgs) error
}

type RenderIndexArgs struct {
	CL       store.Changelog
	Articles []changelog.ParsedArticle
	HasMore  bool
	NextPage int
	PageSize int
	FeedURL  string
}

type RenderArticleListArgs struct {
	CID      store.ChangelogID
	WID      store.WorkspaceID
	Articles []changelog.ParsedArticle
	HasMore  bool
	NextPage int
	PageSize int
}

func New(cfg config.Config) Renderer {
	return &renderer{
		cfg: cfg,
	}
}

type renderer struct {
	cfg config.Config
}

func (r *renderer) RenderArticleList(ctx context.Context, w io.Writer, args RenderArticleListArgs) error {
	articles := parsedArticlesToComponentArticles(args.Articles)

	var nextPageURL string
	if args.HasMore {
		if r.cfg.IsDBMode() {
			nextPageURL = fmt.Sprintf("/?wid=%s&cid=%s&page=%d&page-size=%d", args.WID.String(), args.CID.String(), args.NextPage, args.PageSize)
		} else {
			nextPageURL = fmt.Sprintf("/?page=%d&page-size=%d", args.NextPage, args.PageSize)
		}
	}

	return components.ArticleList(components.ArticleListArgs{
		Articles:    articles,
		NextPageURL: nextPageURL,
	}).Render(ctx, w)
}

func (r *renderer) RenderIndex(ctx context.Context, w io.Writer, args RenderIndexArgs) error {
	var nextPageURL string
	if args.HasMore {
		if r.cfg.IsDBMode() {
			nextPageURL = fmt.Sprintf("/?wid=%s&cid=%s&page=%d&page-size=%d", args.CL.WorkspaceID.String(), args.CL.ID.String(), args.NextPage, args.PageSize)
		} else {
			nextPageURL = fmt.Sprintf("/?page=%d&page-size=%d", args.NextPage, args.PageSize)
		}
	}

	articles := parsedArticlesToComponentArticles(args.Articles)
	return views.Index(views.IndexArgs{
		RSSArgs: components.RSSArgs{
			FeedURL: args.FeedURL,
		},
		MainArgs: layout.MainArgs{
			Title:       args.CL.Title.String,
			Description: args.CL.Subtitle.String,
		},
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
			Articles:    articles,
			NextPageURL: nextPageURL,
		},
	}).Render(ctx, w)
}

func parsedArticlesToComponentArticles(parsed []changelog.ParsedArticle) []components.ArticleArgs {
	articles := make([]components.ArticleArgs, len(parsed))
	for i, a := range parsed {
		buf := new(strings.Builder)
		_, err := io.Copy(buf, a.Content)
		if err != nil {
			continue
		}

		articles[i] = components.ArticleArgs{
			ID:          a.Meta.ID,
			Title:       a.Meta.Title,
			Description: a.Meta.Description,
			PublishedAt: a.Meta.PublishedAt.Format("02 Jan 2006"),
			Content:     buf.String(),
		}
	}

	return articles
}
