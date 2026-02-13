package fullcut

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type Save struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		FullCut *model.FullCut `body:""`
	} `method:"Post"`
}

func (m *Save) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Save) HandlePost(ctx constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	Orm := db.GetDB(ctx)
	//item := &model.FullCut{}
	//err = util.RequestBodyToJSON(context.Request.Body, item)
	//if err != nil {
	//	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	//}
	m.Post.FullCut.OID = m.Organization.ID
	if m.Post.FullCut.ID.IsZero() {
		err = dao.Create(Orm, m.Post.FullCut)
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "添加成功", nil)}, err
	} else {
		err = dao.Save(Orm, m.Post.FullCut)
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
	}
}
