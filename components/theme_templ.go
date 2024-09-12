// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.771
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/jonashiltl/openchangelog/apitypes"

type ThemeArgs struct {
	ColorScheme apitypes.ColorScheme
}

func Theme(args ThemeArgs) templ.Component {
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
		if args.ColorScheme == apitypes.System {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script>\n\t\t\tfunction setTheme(isDark) {\n\t\t\t\tif (isDark) {\n\t\t\t\t\tdocument.querySelector('#theme-container')?.setAttribute('color-scheme', 'dark');\n\t\t\t\t} else {\n\t\t\t\t\tdocument.querySelector('#theme-container')?.setAttribute('color-scheme', 'light');\n\t\t\t\t}\n\t\t\t};\n\n\t\t\t// initialize theme with current scheme\n\t\t\tdocument.addEventListener('DOMContentLoaded', () => \n\t\t\t\tsetTheme(window.matchMedia('(prefers-color-scheme: dark)').matches)\n\t\t\t);\n\n\t\t\t// if loaded throuhg ajax, DOMContentLoaded is not received. Try it manually.\n\t\t\tsetTheme(window.matchMedia('(prefers-color-scheme: dark)').matches)\n\n\t\t\t// listen to theme changes\n\t\t\twindow.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', e => setTheme(e.matches));\n\t\t</script>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"theme-container\" class=\"flex flex-col bg-neutral-50 dark:bg-neutral-950\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if args.ColorScheme != apitypes.System {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" color-scheme=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 string
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(getColorSchemeString(args.ColorScheme))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/theme.templ`, Line: 36, Col: 56}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

func getColorSchemeString(cs apitypes.ColorScheme) string {
	switch cs {
	case apitypes.Light:
		return "light"
	case apitypes.Dark:
		return "dark"
	default:
		return ""
	}
}

var _ = templruntime.GeneratedTemplate
