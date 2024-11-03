package function

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/funcmap"
	"log"
)

type FaviconIco struct {
	Organization  *model.Organization  `mapping:""`
	ContentConfig *model.ContentConfig `mapping:""`
}

func (g *FaviconIco) Call(ctx constrain.IContext) funcmap.IFuncResult {
	/*dir, fileName := filepath.Split(g.ContentConfig.FaviconIco)
	fileNames := strings.Split(fileName, ".")
	if len(fileNames) == 1 && len(fileName) > 0 {
		fileName = fmt.Sprintf("%s%s@64x64", dir, fileName)
	} else if len(fileNames) > 1 {
		fileName = fmt.Sprintf("%s%s@64x64.%s", dir, fileNames[0], fileNames[1])
	}*/
	url, err := oss.ReadUrl(ctx, g.ContentConfig.FaviconIco) //ossurl.CreateUrl(ctx, g.ContentConfig.FaviconIco)
	if err != nil {
		log.Println(err)
	}
	return funcmap.NewStringFuncResult(url)
}
