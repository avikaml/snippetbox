package ui

import (
	"embed"
)
// below is a go directive, not a comment. 
// it stores the files in a filesystem referenced by Files


//go:embed "html" "static"
var Files embed.FS