package field

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type SortField struct {
	Organization *model.Organization `mapping:""`
	Put          struct {
		model.CustomizeField
	} `method:"put"` //更新
}

func (m *SortField) Handle(context constrain.IContext) (constrain.IResult, error) {
	return nil, nil
}
func (m *SortField) HandlePut(ctx constrain.IContext) (constrain.IResult, error) {
	/*customizeFieldGroup := dao.GetBy(db.GetDB(ctx), &model.CustomizeField{}, map[string]any{"GroupID": m.Put.GroupID, "Type": m.Put.Type, "Field": m.Put.Field}).(*model.CustomizeField)
	if customizeFieldGroup.IsZero() == false && customizeFieldGroup.ID != m.Put.ID {
		return nil, errors.New(fmt.Sprintf("字段[%s]已经存在", m.Put.Field))
	}*/
	err := dao.UpdateByPrimaryKey(db.GetDB(ctx), &model.CustomizeField{}, m.Put.ID, map[string]any{
		//"GroupID":  m.Put.GroupID,
		//"Type":     m.Put.Type,
		"Field": m.Put.Field,
		"Sort":  m.Put.Sort,
		//"Extra":    m.Put.Extra,
		//"ParentID": m.Put.ParentID,
	})
	if err != nil {
		return nil, err
	}
	return result.NewSuccess("修改成功"), nil
}
