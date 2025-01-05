package voucher

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"

	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type Index struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		Voucher *model.Voucher `body:""`
	} `method:"Post"`
}

func (m *Index) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Index) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	Orm := db.Orm()
	//item := &model.Voucher{}
	//err = util.RequestBodyToJSON(context.Request.Body, item)
	//if err != nil {
	//	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	//}
	m.Post.Voucher.OID = m.Organization.ID
	err = dao.Create(Orm, m.Post.Voucher)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "添加成功", nil)}, err
}
