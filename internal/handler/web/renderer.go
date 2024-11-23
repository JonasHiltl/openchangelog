package web

import (
	"context"
	"io"
	"strings"

	"github.com/jonashiltl/openchangelog/components"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/handler/web/static"
	"github.com/jonashiltl/openchangelog/internal/handler/web/views"
	"github.com/jonashiltl/openchangelog/internal/handler/web/views/layout"
	"github.com/jonashiltl/openchangelog/internal/parse"
	"github.com/jonashiltl/openchangelog/internal/store"
)

type Renderer interface {
	RenderChangelog(ctx context.Context, w io.Writer, args RenderChangelogArgs) error
	RenderArticleList(ctx context.Context, w io.Writer, args RenderArticleListArgs) error
	RenderWidget(ctx context.Context, w io.Writer, args RenderChangelogArgs) error
	RenderDetails(ctx context.Context, w io.Writer, args RenderDetailsArgs) error
}

type RenderChangelogArgs struct {
	CL           store.Changelog
	ReleaseNotes []parse.ParsedReleaseNote
	HasMore      bool
	CurrentURL   string
	FeedURL      string
	HasMetaKey   bool
}

type RenderArticleListArgs struct {
	CID      store.ChangelogID
	WID      store.WorkspaceID
	Articles []parse.ParsedReleaseNote
	HasMore  bool
	NextPage int
	PageSize int
}

type RenderDetailsArgs struct {
	CL          store.Changelog
	ReleaseNote parse.ParsedReleaseNote
	Prev        parse.ParsedReleaseNote
	Next        parse.ParsedReleaseNote
	FeedURL     string
	HasMetaKey  bool
}

func NewRenderer(cfg config.Config) Renderer {
	return &renderer{
		cfg: cfg,
		css: static.BaseCSS,
	}
}

type renderer struct {
	cfg config.Config
	css string
}

func (r *renderer) RenderArticleList(ctx context.Context, w io.Writer, args RenderArticleListArgs) error {
	articles := parsedArticlesToComponentArticles(args.Articles)
	return components.ArticleList(components.ArticleListArgs{
		Articles: articles,
	}).Render(ctx, w)
}

func (r *renderer) RenderChangelog(ctx context.Context, w io.Writer, args RenderChangelogArgs) error {
	notes := parsedArticlesToComponentArticles(args.ReleaseNotes)
	return views.Index(views.IndexArgs{
		RSSArgs: components.RSSArgs{
			FeedURL: args.FeedURL,
		},
		SearchButtonArgs: components.SearchButtonArgs{
			Show:       args.CL.Searchable,
			HasMetaKey: args.HasMetaKey,
		},
		ChangelogContainerArgs: components.ChangelogContainerArgs{
			CurrentURL:     args.CurrentURL,
			HasMoreArticle: args.HasMore,
		},
		MainArgs: layout.MainArgs{
			Title:       args.CL.Title.V(),
			Description: args.CL.Subtitle.V(),
			CSS:         r.css,
		},
		ThemeArgs: components.ThemeArgs{
			ColorScheme: args.CL.ColorScheme.ToApiTypes(),
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
			Articles: notes,
		},
		FooterArgs: components.FooterArgs{
			HidePoweredBy: args.CL.HidePoweredBy,
		},
	}).Render(ctx, w)
}

func (r *renderer) RenderWidget(ctx context.Context, w io.Writer, args RenderChangelogArgs) error {
	articles := parsedArticlesToComponentArticles(args.ReleaseNotes)
	return views.Widget(views.WidgetArgs{
		CSS: r.css,
		ChangelogContainerArgs: components.ChangelogContainerArgs{
			CurrentURL:     args.CurrentURL,
			HasMoreArticle: args.HasMore,
		},
		HeaderArgs: components.HeaderArgs{
			Title:    args.CL.Title,
			Subtitle: args.CL.Subtitle,
		},
		ArticleListArgs: components.ArticleListArgs{
			Articles: articles,
		},
		FooterArgs: components.FooterArgs{
			HidePoweredBy: args.CL.HidePoweredBy,
		},
	}).Render(ctx, w)
}

func (r *renderer) RenderDetails(ctx context.Context, w io.Writer, args RenderDetailsArgs) error {
	articles := parsedArticlesToComponentArticles([]parse.ParsedReleaseNote{
		args.ReleaseNote, args.Prev, args.Next,
	})
	return views.Details(views.DetailsArgs{
		RSSArgs: components.RSSArgs{
			FeedURL: args.FeedURL,
		},
		SearchButtonArgs: components.SearchButtonArgs{
			Show:       args.CL.Searchable,
			HasMetaKey: args.HasMetaKey,
		},
		MainArgs: layout.MainArgs{
			Title:       args.CL.Title.V(),
			Description: args.CL.Subtitle.V(),
			CSS:         r.css,
		},
		HeaderArgs: components.HeaderArgs{
			Title:    args.CL.Title,
			Subtitle: args.CL.Subtitle,
			ShowBack: true,
		},
		ThemeArgs: components.ThemeArgs{
			ColorScheme: args.CL.ColorScheme.ToApiTypes(),
		},
		Logo: components.Logo{
			Src:    args.CL.LogoSrc,
			Width:  args.CL.LogoWidth,
			Height: args.CL.LogoHeight,
			Alt:    args.CL.LogoAlt,
			Link:   args.CL.LogoLink,
		},
		ArticleArgs: articles[0],
		Prev:        articles[1],
		Next:        articles[2],
		FooterArgs: components.FooterArgs{
			HidePoweredBy: args.CL.HidePoweredBy,
		},
	}).Render(ctx, w)
}

func parsedArticlesToComponentArticles(parsed []parse.ParsedReleaseNote) []components.ArticleArgs {
	articles := make([]components.ArticleArgs, len(parsed))
	for i, a := range parsed {
		article := components.ArticleArgs{
			ID:          a.Meta.ID,
			Title:       a.Meta.Title,
			Description: a.Meta.Description,
			PublishedAt: a.Meta.PublishedAt,
			Tags:        a.Meta.Tags,
		}
		if a.Content != nil {
			buf := new(strings.Builder)
			_, err := io.Copy(buf, a.Content)
			if err != nil {
				continue
			}
			article.Content = buf.String()
		}

		articles[i] = article
	}

	return articles
}
