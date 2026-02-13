package express

import (
	"fmt"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type DeliveryList struct {
	Post struct {
		Name string
	} `method:"post"`
}

func (m *DeliveryList) Handle(ctx constrain.IContext) (r constrain.IResult, err error) {

	list := dao.Find(db.GetDB(ctx), entity.ExpressCompany).List()

	return result.NewData(list), nil
}
func (m *DeliveryList) HandlePost(ctx constrain.IContext) (r constrain.IResult, err error) {

	d := dao.Find(db.GetDB(ctx), entity.ExpressCompany)
	if len(m.Post.Name) > 0 {
		d.Where(`"Name" like ?`, fmt.Sprintf("%%%s%%", m.Post.Name))
	}
	list := d.List()
	return result.NewData(list), nil
}
