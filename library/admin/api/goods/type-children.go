package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type TypeChildren struct {
	OIDMapping *entity.SessionMappingData `mapping:""`
	Get        struct {
		ParentID uint `form:"ParentID"`
	} `method:"get"`
}

func (m *TypeChildren) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	list := dao.Find(db.Orm(), &model.GoodsType{}).Where(`"OID"=? and "ParentID"=?`, m.OIDMapping.OID, m.Get.ParentID).List()

	return result.NewData(map[string]any{"List": list}), nil
}
