package static

import (
	_ "embed"
)

//go:embed base.css
var BaseCSS string

//go:embed admin.css
var AdminCSS string
