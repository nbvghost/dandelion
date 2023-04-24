package goods

import (
	"fmt"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/internal/repository"
	"github.com/nbvghost/gpa/types"
	"github.com/pkg/errors"
	"strings"
)

type AttributesService struct {
}

func (service AttributesService) AllAttributesName() ([]*extends.GoodsAttributesNameInfo, error) {

	return repository.GoodsAttributes.QueryGoodsAttributesNameInfo()
}
func (service AttributesService) AllAttributesByName(name string) ([]*extends.GoodsAttributesValueInfo, error) {

	return repository.GoodsAttributes.QueryGoodsAttributesValueInfoByName(name)
}
func (service AttributesService) DeleteGoodsAttributes(ID types.PrimaryKey) error {

	return repository.GoodsAttributes.DeleteByID(ID).Err
}

func (service AttributesService) AddGoodsAttributes(goodsID, groupID types.PrimaryKey, name, value string) error {
	if goodsID == 0 || groupID == 0 {
		return errors.New(fmt.Sprintf("产品ID不能为空或组ID不能为空"))
	}
	if strings.EqualFold(name, "") || strings.EqualFold(value, "") {
		return nil
	}
	hasAttr, err := repository.GoodsAttributes.GetByGoodsIDAndName(goodsID, name)
	if hasAttr.IsZero() == false {
		return errors.New(fmt.Sprintf("属性名：%v已经存在", name))
	}
	err = repository.GoodsAttributes.Save(&model.GoodsAttributes{
		GoodsID: goodsID,
		GroupID: groupID,
		Name:    name,
		Value:   value,
	})
	if err != nil {
		return err
	}
	return nil
}
func (service AttributesService) ListGoodsAttributesGroupByGoodsID(goodsID types.PrimaryKey) ([]*model.GoodsAttributesGroup, error) {

	return repository.GoodsAttributesGroup.FindByGoodsID(goodsID)
}
func (service AttributesService) GetGoodsAttributesGroup(ID types.PrimaryKey) types.IEntity {
	return repository.GoodsAttributesGroup.GetByID(ID)
}
func (service AttributesService) DeleteGoodsAttributesGroup(ID types.PrimaryKey) error {
	attrs, err := service.ListGoodsAttributesByGroupID(ID)
	if err != nil {
		return err
	}
	if len(attrs) > 0 {
		return errors.New(fmt.Sprintf("属性组包含子属性，无法删除"))
	}
	del := repository.GoodsAttributesGroup.DeleteByID(ID)
	return del.Err
}
func (service AttributesService) ListGoodsAttributesByGroupID(attributesGroupID types.PrimaryKey) ([]*model.GoodsAttributes, error) {
	return repository.GoodsAttributes.FindByGroupID(attributesGroupID)
}
func (service AttributesService) ChangeGoodsAttributesGroup(id types.PrimaryKey, groupName string) error {
	if id == 0 {
		return errors.New(fmt.Sprintf("ID不能为空"))
	}
	if strings.EqualFold(groupName, "") {
		return nil
	}
	hasAttr, err := repository.GoodsAttributesGroup.GetByName(groupName)
	if err != nil {
		return err
	}
	if hasAttr.IsZero() == false {
		return errors.New(fmt.Sprintf("属性名：%v已经存在", groupName))
	}

	update := repository.GoodsAttributesGroup.UpdateByID(id, map[string]interface{}{"Name": groupName})
	if update.Err != nil {
		return err
	}
	return nil
}
func (service AttributesService) AddGoodsAttributesGroup(goodsID types.PrimaryKey, groupName string) error {
	if goodsID == 0 {
		return errors.New(fmt.Sprintf("产品ID不能为空"))
	}
	if strings.EqualFold(groupName, "") {
		return nil
	}
	hasAttr, err := repository.GoodsAttributesGroup.GetByGoodsIDAndName(goodsID, groupName)
	if err != nil {
		return err
	}
	if hasAttr.IsZero() == false {
		return errors.New(fmt.Sprintf("属性名：%v已经存在", groupName))
	}
	err = repository.GoodsAttributesGroup.Save(&model.GoodsAttributesGroup{
		GoodsID: goodsID,
		Name:    groupName,
	})
	if err != nil {
		return err
	}
	return nil
}
