package goods

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
	"strings"
)

type ChangeGoodsType struct {
	OIDMapping entity.SessionMappingData `mapping:""`
	Post       struct {
		*model.GoodsType
	} `method:"post"`
}

func (g *ChangeGoodsType) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	return nil, err
}
func (g *ChangeGoodsType) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//item := &model.GoodsType{}
	//err = util.RequestBodyToJSON(context.Request.Body, item)
	//log.Println(err)

	if strings.HasPrefix(g.Post.GoodsType.Badge, oss.TempFilePrefix) {
		fileBytes, err := oss.GetTempFile(g.Post.GoodsType.Badge)
		if err != nil {
			return nil, err
		}
		file, err := oss.UploadFile(context, fileBytes, fmt.Sprintf("goodstype/%d", g.Post.ID), "", true, "badge")
		if err != nil {

			return nil, err
		}
		g.Post.GoodsType.Badge = file.Data.Path
	}
	if strings.HasPrefix(g.Post.GoodsType.Image, oss.TempFilePrefix) {
		fileBytes, err := oss.GetTempFile(g.Post.GoodsType.Image)
		if err != nil {
			return nil, err
		}
		file, err := oss.UploadFile(context, fileBytes, fmt.Sprintf("goodstype/%d", g.Post.ID), "", true, "image")
		if err != nil {

			return nil, err
		}
		g.Post.GoodsType.Image = file.Data.Path
	}

	err = service.Goods.GoodsType.ChangeGoodsType(g.OIDMapping.OID, g.Post.GoodsType)
	if err != nil {
		return nil, err
	}

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}
