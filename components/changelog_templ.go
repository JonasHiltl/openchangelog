// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.771
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "html/template"

type ChangelogContainerArgs struct {
	CurrentURL     string
	HasMoreArticle bool
}

// Contains the article list and footer
func ChangelogContainer(args ChangelogContainerArgs) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<main id=\"changelog-container\" class=\"o-mx-4 sm:o-mx-0\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"skeleton\" class=\"o-hidden\"><div class=\"o-animate-pulse o-w-full o-space-y-4 o-mt-12\"><div class=\"o-w-3/4 o-h-9 o-rounded o-bg-black/10 dark:o-bg-white/10\"></div><div class=\"o-w-full o-h-5 o-rounded o-bg-black/10 dark:o-bg-white/10\"></div><div class=\"o-w-full o-h-32 o-rounded o-bg-black/10 dark:o-bg-white/10\"></div></div></div></main>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if args.HasMoreArticle {
			templ_7745c5c3_Err = templ.FromGoHTML(infiniteScrollTemplate, args.CurrentURL).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		return templ_7745c5c3_Err
	})
}

var infiniteScrollTemplate = template.Must(template.New("infiniteScrollTemplate").Parse(`
<script>
	const skeleton = document.getElementById("skeleton")
	const container = document.getElementById("changelog-container")
	const currentPageURL = new URL("{{ . }}")
	const params = currentPageURL.searchParams;
	let isLoadingNextPage = false
	let hasMore = true

	function loadStarted() {
		isLoadingNextPage = true
		skeleton.style.display = "block"
	}

	function loadEnded() {
		isLoadingNextPage = false
		skeleton.style.display = "none"
	}

	async function loadNextPage(url) {
		try {
			loadStarted()
			const currentPage = parseInt(params.get('page')) || 1;
			params.set('page', currentPage + 1);
 			// don't load full html page, only articles list
			params.set('articles', "true");

			const res = await fetch(currentPageURL)
			if (!res.ok || res.status === 204) {
				hasMore = false
			}
			const newArticles = await res.text()
			skeleton.insertAdjacentHTML('beforebegin', newArticles)
		} finally {
			loadEnded()
		}
	}

	window.addEventListener("scroll", () => {
		const bufferPx = 100 // start loading before reaching end of container
		const endReached = window.scrollY + window.innerHeight + bufferPx >= container.scrollHeight
		if(endReached && !isLoadingNextPage && hasMore){
			loadNextPage()
		}
	})
</script>
`,
))

var _ = templruntime.GeneratedTemplate
