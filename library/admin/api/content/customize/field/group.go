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

type Group struct {
	Organization *model.Organization `mapping:""`

	Post struct {
		Name string
	} `method:"post"` //添加
	Put struct {
		ID   dao.PrimaryKey
		Name string
	} `method:"put"` //更新
	Delete struct {
		ID dao.PrimaryKey
	} `method:"delete"` //删除
}

func (m *Group) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	list := dao.Find(db.Orm(), &model.CustomizeFieldGroup{}).Where(`"OID"=?`, m.Organization.ID).List()
	return result.NewData(map[string]any{"List": list}), nil
}
func (m *Group) HandlePost(context constrain.IContext) (constrain.IResult, error) {
	customizeFieldGroup := dao.GetBy(db.Orm(), &model.CustomizeFieldGroup{}, map[string]any{"OID": m.Organization.ID, "Name": m.Post.Name}).(*model.CustomizeFieldGroup)
	if !customizeFieldGroup.IsZero() {
		return nil, errors.New(fmt.Sprintf("分组名[%s]已经存在", m.Post.Name))
	}

	err := dao.Create(db.Orm(), &model.CustomizeFieldGroup{
		OID:  m.Organization.ID,
		Name: m.Post.Name,
	})
	if err != nil {
		return nil, err
	}
	return result.NewSuccess("添加成功"), nil
}
func (m *Group) HandlePut(context constrain.IContext) (constrain.IResult, error) {
	customizeFieldGroup := dao.GetBy(db.Orm(), &model.CustomizeFieldGroup{}, map[string]any{"OID": m.Organization.ID, "Name": m.Put.Name}).(*model.CustomizeFieldGroup)
	if customizeFieldGroup.IsZero() == false && customizeFieldGroup.ID != m.Put.ID {
		return nil, errors.New(fmt.Sprintf("分组名[%s]已经存在", m.Put.Name))
	}
	err := dao.UpdateByPrimaryKey(db.Orm(), &model.CustomizeFieldGroup{}, m.Put.ID, map[string]any{"Name": m.Put.Name})
	if err != nil {
		return nil, err
	}
	return result.NewSuccess("修改成功"), nil
}
func (m *Group) HandleDelete(context constrain.IContext) (constrain.IResult, error) {
	has := dao.GetBy(db.Orm(), &model.CustomizeField{}, map[string]any{"GroupID": m.Delete.ID})
	if !has.IsZero() {
		return nil, errors.New(fmt.Sprintf("分组包含子项内容，无法删除"))
	}
	err := dao.DeleteByPrimaryKey(db.Orm(), &model.CustomizeFieldGroup{}, m.Delete.ID)
	if err != nil {
		return nil, err
	}
	return result.NewSuccess("删除成功"), nil
}
