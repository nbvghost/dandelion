package goods

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/gpa/types"
	"github.com/nbvghost/gweb"
	"github.com/nbvghost/tool/object"
)

type SpecificationService struct {
}

func (service SpecificationService) GetSpecification(ID types.PrimaryKey) *model.Specification {
	Orm := singleton.Orm()
	return dao.GetByPrimaryKey(Orm, &model.Specification{}, ID).(*model.Specification)
}
func (service SpecificationService) DeleteSpecification(ID types.PrimaryKey) error {
	Orm := singleton.Orm()
	err := dao.DeleteByPrimaryKey(Orm, &model.Specification{}, ID)
	return err
}
func (service SpecificationService) ChangeSpecification(context *gweb.Context) (r gweb.Result, err error) {
	Orm := singleton.Orm()
	//GoodsID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)
	GoodsID := object.ParseUint(context.PathParams["GoodsID"])
	item := &model.Specification{}
	err = util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	}
	err = dao.UpdateByPrimaryKey(Orm, &model.Specification{}, types.PrimaryKey(GoodsID), item)
	return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}
