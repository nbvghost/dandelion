package webpicture

import _ "embed"

//go:embed mac/cwebp
var cwebpBytes []byte

//go:embed mac/gif2webp
var gif2webpBytes []byte
