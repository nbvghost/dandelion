//go:build aix || dragonfly || freebsd || wasip1 || linux || netbsd || openbsd || solaris

package webpicture

import _ "embed"

//go:embed linux/cwebp
var cwebpBytes []byte

//go:embed linux/gif2webp
var gif2webpBytes []byte
