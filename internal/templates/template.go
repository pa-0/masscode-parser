package templates

import (
	"embed"
)

//go:embed html/template.html
var TemplateFS embed.FS
