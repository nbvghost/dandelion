package embed

import (
	_ "embed"
)

//go:embed template/RestPasswordEmailTemplate.gohtml
var RestPasswordEmailTemplate string
