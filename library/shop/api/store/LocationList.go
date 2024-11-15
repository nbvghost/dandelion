package store

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type LocationList struct {
	User *model.User `mapping:""`
	Get  struct {
		Latitude  float64 `uri:"Latitude"`
		Longitude float64 `uri:"Longitude"`
	} `method:"Get"`
}

func (m *LocationList) Handle(context constrain.IContext) (r constrain.IResult, err error) {

	list := service.Company.Store.LocationList(m.Get.Latitude, m.Get.Longitude)

	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: list}}, nil
}
