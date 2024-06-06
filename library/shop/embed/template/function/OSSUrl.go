package function

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/library/funcmap"
	"github.com/nbvghost/tool/object"
	"log"
	"strings"
)

type OSSUrl struct {
	Path any `arg:""`
}

func (g *OSSUrl) Call(ctx constrain.IContext) funcmap.IFuncResult {
	path := object.ParseString(g.Path)
	if strings.EqualFold(path, "") {
		path = "/default"
	}
	url, err := oss.ReadUrl(ctx, path) //ossurl.CreateUrl(ctx, path)
	if err != nil {
		log.Println(err)
	}
	return funcmap.NewStringFuncResult(url)
}
