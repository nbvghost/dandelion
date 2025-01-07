package goods

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
	"strings"
)

type AddGoodsType struct {
	OIDMapping *entity.SessionMappingData `mapping:""`
	Get        struct {
		*model.GoodsType
	} `method:"Get"`
	Post struct {
		*model.GoodsType
	} `method:"Post"`
}

func (m *AddGoodsType) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	return nil, err
}
func (m *AddGoodsType) HandlePost(context constrain.IContext) (constrain.IResult, error) {
	gt, err := service.Goods.GoodsType.AddGoodsType(m.OIDMapping.OID, m.Post.GoodsType)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(m.Post.GoodsType.Badge, oss.TempFilePrefix) {
		fileBytes, err := oss.GetTempFile(m.Post.GoodsType.Badge)
		if err != nil {
			return nil, err
		}
		file, err := oss.UploadFile(context, fileBytes, fmt.Sprintf("goodstype/%d", gt.ID), "", true, "badge")
		if err != nil {

			return nil, err
		}
		m.Post.GoodsType.Badge = file.Data.Path
	}
	if strings.HasPrefix(m.Post.GoodsType.Image, oss.TempFilePrefix) {
		fileBytes, err := oss.GetTempFile(m.Post.GoodsType.Image)
		if err != nil {
			return nil, err
		}
		file, err := oss.UploadFile(context, fileBytes, fmt.Sprintf("goodstype/%d", gt.ID), "", true, "image")
		if err != nil {

			return nil, err
		}
		m.Post.GoodsType.Image = file.Data.Path
	}
	err = dao.UpdateByPrimaryKey(db.Orm(), &model.GoodsType{}, gt.ID, map[string]any{"Badge": m.Post.GoodsType.Badge, "Image": m.Post.GoodsType.Image})
	if err != nil {
		return nil, err
	}
	return result.NewData(gt), err
}
