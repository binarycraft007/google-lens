package main

import (
	_ "image/jpeg"
	_ "image/png"
	"strings"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/binarycraft007/lens"
)

func main() {
	lensApi := lens.NewLens()
	bytes, err := lensApi.Scan(strings.NewReader(data))
	if err != nil {
		panic(err)
	}
	_ = bytes
}

const data = `
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
`
