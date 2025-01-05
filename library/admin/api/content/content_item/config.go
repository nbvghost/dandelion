package content_item

import (
	"encoding/json"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type Config struct {
	Organization *model.Organization `mapping:""`
	POST         struct {
		ContentItemID dao.PrimaryKey
		Type          string
		TemplateName  string
	} `method:"POST"`
}

func (m *Config) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Config) HandlePost(context constrain.IContext) (constrain.IResult, error) {
	c := dao.GetByPrimaryKey(db.Orm(), entity.ContentItem, m.POST.ContentItemID).(*model.ContentItem)
	config, err := json.Marshal(map[string]any{"Type": m.POST.Type, "TemplateName": m.POST.TemplateName})
	if err != nil {
		return nil, err
	}
	c.Config = string(config)

	err = dao.UpdateByPrimaryKey(db.Orm(), entity.ContentItem, c.ID, map[string]any{"Config": c.Config})
	if err != nil {
		return nil, err
	}
	return result.NewSuccess("修改成功"), nil
}
