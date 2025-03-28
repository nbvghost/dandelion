package express

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type DeliveryList struct {
	Admin *model.Admin `mapping:""`
	Post  struct {
		Name string
	} `method:"post"`
}

func (m *DeliveryList) Handle(context constrain.IContext) (r constrain.IResult, err error) {

	list := dao.Find(db.Orm(), entity.ExpressCompany).List()

	return result.NewData(list), nil
}
func (m *DeliveryList) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {

	d := dao.Find(db.Orm(), entity.ExpressCompany)
	if len(m.Post.Name) > 0 {
		d.Where(`"Name" like ?`, fmt.Sprintf("%%%s%%", m.Post.Name))
	}
	list := d.List()
	return result.NewData(list), nil
}
