package field

import (
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type Field struct {
	Organization *model.Organization `mapping:""`
	Get          struct {
		model.CustomizeField
	} `method:"get"` //添加
	Post struct {
		model.CustomizeField
	} `method:"post"` //添加
	Put struct {
		model.CustomizeField
	} `method:"put"` //更新
	Delete struct {
		ID dao.PrimaryKey
	} `method:"delete"` //删除
}

func (m *Field) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	list := dao.Find(db.Orm(), &model.CustomizeField{}).Where(map[string]any{"OID": m.Organization.ID}).Order(`"Sort","UpdatedAt" desc`).List()
	return result.NewData(map[string]any{"List": list}), nil
}
func (m *Field) HandlePost(context constrain.IContext) (constrain.IResult, error) {
	/*customizeFieldGroup := dao.GetBy(db.Orm(), &model.CustomizeField{}, map[string]any{"GroupID": m.Post.GroupID, "Type": m.Post.Type, "Field": m.Post.Field}).(*model.CustomizeField)
	if !customizeFieldGroup.IsZero() {
		return nil, errors.New(fmt.Sprintf("字段[%s]已经存在", m.Post.Field))
	}*/

	err := dao.Create(db.Orm(), &model.CustomizeField{
		OID:      m.Organization.ID,
		GroupID:  m.Post.GroupID,
		Type:     m.Post.Type,
		Field:    m.Post.Field,
		Extra:    m.Post.Extra,
		ParentID: m.Post.ParentID,
	})
	if err != nil {
		return nil, err
	}
	return result.NewSuccess("添加成功"), nil
}
func (m *Field) HandlePut(context constrain.IContext) (constrain.IResult, error) {
	/*customizeFieldGroup := dao.GetBy(db.Orm(), &model.CustomizeField{}, map[string]any{"GroupID": m.Put.GroupID, "Type": m.Put.Type, "Field": m.Put.Field}).(*model.CustomizeField)
	if customizeFieldGroup.IsZero() == false && customizeFieldGroup.ID != m.Put.ID {
		return nil, errors.New(fmt.Sprintf("字段[%s]已经存在", m.Put.Field))
	}*/
	err := dao.UpdateByPrimaryKey(db.Orm(), &model.CustomizeField{}, m.Put.ID, map[string]any{
		//"GroupID":  m.Put.GroupID,
		//"Type":     m.Put.Type,
		"Field": m.Put.Field,
		//"Extra":    m.Put.Extra,
		//"ParentID": m.Put.ParentID,
	})
	if err != nil {
		return nil, err
	}
	return result.NewSuccess("修改成功"), nil
}
func (m *Field) HandleDelete(context constrain.IContext) (constrain.IResult, error) {
	field := dao.GetByPrimaryKey(db.Orm(), &model.CustomizeField{}, m.Delete.ID).(*model.CustomizeField)
	if field.ParentID == 0 {
		hasList := dao.Find(db.Orm(), &model.CustomizeField{}).Where(`"ParentID"=?`, m.Delete.ID).List()
		if len(hasList) > 0 {
			return nil, errors.New(fmt.Sprintf("分组包含子项内容，无法删除"))
		}
	}
	err := dao.DeleteByPrimaryKey(db.Orm(), &model.CustomizeField{}, m.Delete.ID)
	if err != nil {
		return nil, err
	}
	return result.NewSuccess("删除成功"), nil
}
