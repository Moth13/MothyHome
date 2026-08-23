package web

import "embed"

//go:embed templates/index.html
var templateFS embed.FS

//go:embed static/styles.css
var staticFS embed.FS
