package content

import (
	"fmt"
	"time"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
	"gorm.io/gorm"
)

type Like struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		ID dao.PrimaryKey `form:"ID"`
	} `method:"Post"`
}

func (m *Like) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Like) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	contextValue := contexext.FromContext(ctx)
	_, err := ctx.Redis().Get(ctx, fmt.Sprintf("content:like:%d:%s", m.Post.ID, util.GetIP(contextValue.Request)))
	if err == nil {
		//说明已经点赞
		return &result.JsonResult{Data: result.ActionResult{Code: result.Fail, Message: ""}}, nil
	}
	now := time.Now()
	d := now.Add(time.Hour * 24).Sub(now)
	err = ctx.Redis().Set(ctx, fmt.Sprintf("content:like:%d:%s", m.Post.ID, util.GetIP(contextValue.Request)), time.Now().Unix(), d)
	if err != nil {
		return nil, err
	}
	err = db.GetDB(ctx).Model(model.Content{}).Where(map[string]any{"ID": m.Post.ID}).Updates(map[string]any{"CountLike": gorm.Expr(`"CountLike"+1`)}).Error
	if err != nil {
		return nil, err
	}
	return &result.JsonResult{Data: result.ActionResult{}}, nil
}
