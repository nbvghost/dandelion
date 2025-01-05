package timesell

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type Hash struct {
	Organization *model.Organization `mapping:""`
	Get          struct {
		Hash string `uri:"Hash"`
	} `method:"Get"`
}

func (m *Hash) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Hash) HandleGet(context constrain.IContext) (r constrain.IResult, err error) {
	//Orm := db.Orm()
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//Hash := context.PathParams["Hash"]
	item := service.Activity.TimeSell.GetTimeSellByHash(m.Get.Hash, m.Organization.ID)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", item)}, err
}
