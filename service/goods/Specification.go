package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type SpecificationService struct {
}

func (service SpecificationService) GetSpecification(ID dao.PrimaryKey) *model.Specification {
	Orm := db.Orm()
	return dao.GetByPrimaryKey(Orm, &model.Specification{}, ID).(*model.Specification)
}
func (service SpecificationService) DeleteSpecification(ID dao.PrimaryKey) error {
	Orm := db.Orm()
	err := dao.DeleteByPrimaryKey(Orm, &model.Specification{}, ID)
	return err
}
func (service SpecificationService) ChangeSpecification(context constrain.IContext, goodsID dao.PrimaryKey) (r constrain.IResult, err error) {
	//Orm := db.Orm()
	//GoodsID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)
	//GoodsID := object.ParseUint(context.PathParams["GoodsID"])
	//item := &model.Specification{}
	/*err = util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	}
	err = dao.UpdateByPrimaryKey(Orm, &model.Specification{}, dao.PrimaryKey(GoodsID), item)
	*/
	//return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
	return result.NewSuccess("修改成功"), err
}
