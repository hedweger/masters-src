package templates

import "embed"

//go:embed data/*.tmpl
var templateFiles embed.FS

