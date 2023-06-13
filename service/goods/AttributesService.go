package goods

import (
	"fmt"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/gpa"
	"github.com/nbvghost/gpa/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type AttributesService struct {
}

////FindByGoodsID       func(goodsID types.PrimaryKey) ([]*model.GoodsAttributes, error)            `gpa:"AutoCrate"`
//FindByGroupID       func(groupID types.PrimaryKey) ([]*model.GoodsAttributes, error)            `gpa:"AutoCrate"`
//GetByGoodsIDAndName func(goodsID types.PrimaryKey, name string) (*model.GoodsAttributes, error) `gpa:"AutoCrate"`

//FindByGoodsID       func(goodsID types.PrimaryKey) ([]*model.GoodsAttributesGroup, error)            `gpa:"AutoCrate"`
//GetByGoodsIDAndName func(goodsID types.PrimaryKey, name string) (*model.GoodsAttributesGroup, error) `gpa:"AutoCrate"`
//GetByName           func(name string) (*model.GoodsAttributesGroup, error)                           `gpa:"AutoCrate"`

func (service AttributesService) FindGroupByGoodsID(Orm *gorm.DB, goodsID types.PrimaryKey) []*model.GoodsAttributesGroup {
	list := make([]*model.GoodsAttributesGroup, 0)
	Orm.Where(`"GoodsID"=?`, goodsID).Find(&list) //SelectOne(user, "select * from User where Tel=?", Tel)
	return list
}
func (service AttributesService) GetGroupByName(Orm *gorm.DB, name string) *model.GoodsAttributesGroup {
	item := &model.GoodsAttributesGroup{}
	Orm.Where(`"Name"=?`, name).First(item) //SelectOne(user, "select * from User where Tel=?", Tel)
	return item
}
func (service AttributesService) GetGroupByGoodsIDAndName(Orm *gorm.DB, goodsID types.PrimaryKey, name string) *model.GoodsAttributesGroup {
	item := &model.GoodsAttributesGroup{}
	Orm.Where(`"GoodsID"=? and "Name"=?`, goodsID, name).First(item) //SelectOne(user, "select * from User where Tel=?", Tel)
	return item
}

func (service AttributesService) FindByGoodsID(Orm *gorm.DB, goodsID types.PrimaryKey) []*model.GoodsAttributes {
	list := make([]*model.GoodsAttributes, 0)
	Orm.Where(`"GoodsID"=?`, goodsID).Find(&list) //SelectOne(user, "select * from User where Tel=?", Tel)
	return list
}
func (service AttributesService) FindByGroupID(Orm *gorm.DB, groupID types.PrimaryKey) []*model.GoodsAttributes {
	list := make([]*model.GoodsAttributes, 0)
	Orm.Where(`"GroupID"=?`, groupID).Find(&list) //SelectOne(user, "select * from User where Tel=?", Tel)
	return list
}
func (service AttributesService) GetByGoodsIDAndName(Orm *gorm.DB, goodsID types.PrimaryKey, name string) *model.GoodsAttributes {
	item := &model.GoodsAttributes{}
	Orm.Where(`"GoodsID"=? and "Name"=?`, goodsID, name).First(item) //SelectOne(user, "select * from User where Tel=?", Tel)
	return item
}

func (service AttributesService) QueryGoodsAttributesNameInfo() ([]*extends.GoodsAttributesNameInfo, error) {
	rows, err := db.Orm().Raw(`select * from (select "Name",count("Name") as "Num" from "GoodsAttributes" group by "Name") as m order by m."Num" desc`, nil).Rows()
	if err != nil {
		return nil, err
	}
	d, err := gpa.ScanRows(rows, reflect.TypeOf(new(extends.GoodsAttributesNameInfo)), true)
	if err != nil {
		return nil, err
	}
	//list := gpa.Rows("select * from (select Value,count(Value) as Num from GoodsAttributes where Name=? group by Value) as m order by m.Num desc", []interface{}{name}, &extends.GoodsAttributesValueInfo{})
	return d.([]*extends.GoodsAttributesNameInfo), err

	//list := gpa.Rows("select * from (select Name,count(Name) as Num from GoodsAttributes group by Name) as m order by m.Num desc", nil, &extends.GoodsAttributesNameInfo{})
	//return list.([]*extends.GoodsAttributesNameInfo)
}
func (service AttributesService) QueryGoodsAttributesValueInfoByName(name string) ([]*extends.GoodsAttributesValueInfo, error) {
	rows, err := db.Orm().Raw(`select * from (select "Value",count("Value") as "Num" from "GoodsAttributes" where "Name"=? group by "Value") as m order by m."Num" desc`, name).Rows()
	if err != nil {
		return nil, err
	}
	d, err := gpa.ScanRows(rows, reflect.TypeOf(new(extends.GoodsAttributesValueInfo)), true)
	if err != nil {
		return nil, err
	}
	//list := gpa.Rows("select * from (select Value,count(Value) as Num from GoodsAttributes where Name=? group by Value) as m order by m.Num desc", []interface{}{name}, &extends.GoodsAttributesValueInfo{})
	return d.([]*extends.GoodsAttributesValueInfo), err
}

func (service AttributesService) AllAttributesName() ([]*extends.GoodsAttributesNameInfo, error) {

	return service.QueryGoodsAttributesNameInfo()
}
func (service AttributesService) AllAttributesByName(name string) ([]*extends.GoodsAttributesValueInfo, error) {

	return service.QueryGoodsAttributesValueInfoByName(name)
}
func (service AttributesService) DeleteGoodsAttributes(ID types.PrimaryKey) error {

	return dao.DeleteByPrimaryKey(db.Orm(), &model.GoodsAttributes{}, ID) //repository.GoodsAttributes.DeleteByID(ID).Err
}

func (service AttributesService) AddGoodsAttributes(goodsID, groupID types.PrimaryKey, name, value string) error {
	if goodsID == 0 || groupID == 0 {
		return errors.New(fmt.Sprintf("产品ID不能为空或组ID不能为空"))
	}
	if strings.EqualFold(name, "") || strings.EqualFold(value, "") {
		return nil
	}
	hasAttr := service.GetByGoodsIDAndName(db.Orm(), goodsID, name) //repository.GoodsAttributes.GetByGoodsIDAndName(goodsID, name)
	if hasAttr.IsZero() == false {
		return errors.New(fmt.Sprintf("属性名：%v已经存在", name))
	}
	err := dao.Create(db.Orm(), &model.GoodsAttributes{
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
func (service AttributesService) ListGoodsAttributesGroupByGoodsID(goodsID types.PrimaryKey) []*model.GoodsAttributesGroup {

	return service.FindGroupByGoodsID(db.Orm(), goodsID) //repository.GoodsAttributesGroup.FindByGoodsID(goodsID)
}
func (service AttributesService) GetGoodsAttributesGroup(ID types.PrimaryKey) types.IEntity {
	return dao.GetByPrimaryKey(db.Orm(), &model.GoodsAttributesGroup{}, ID) //repository.GoodsAttributesGroup.GetByID(ID)
}
func (service AttributesService) DeleteGoodsAttributesGroup(ID types.PrimaryKey) error {
	attrs := service.ListGoodsAttributesByGroupID(ID)

	if len(attrs) > 0 {
		return errors.New(fmt.Sprintf("属性组包含子属性，无法删除"))
	}
	del := dao.DeleteByPrimaryKey(db.Orm(), &model.GoodsAttributesGroup{}, ID) //repository.GoodsAttributesGroup.DeleteByID(ID)
	return del
}
func (service AttributesService) ListGoodsAttributesByGroupID(attributesGroupID types.PrimaryKey) []*model.GoodsAttributes {
	return service.FindByGroupID(db.Orm(), attributesGroupID) //repository.GoodsAttributes.FindByGroupID(attributesGroupID)
}
func (service AttributesService) ChangeGoodsAttributesGroup(id types.PrimaryKey, groupName string) error {
	if id == 0 {
		return errors.New(fmt.Sprintf("ID不能为空"))
	}
	if strings.EqualFold(groupName, "") {
		return nil
	}
	hasAttr := service.GetGroupByName(db.Orm(), groupName) //repository.GoodsAttributesGroup.GetByName(groupName)
	if hasAttr.IsZero() == false {
		return errors.New(fmt.Sprintf("属性名：%v已经存在", groupName))
	}

	err := dao.UpdateByPrimaryKey(db.Orm(), &model.GoodsAttributesGroup{}, id, map[string]interface{}{"Name": groupName}) //repository.GoodsAttributesGroup.UpdateByID(id, map[string]interface{}{"Name": groupName})

	return err
}
func (service AttributesService) AddGoodsAttributesGroup(goodsID types.PrimaryKey, groupName string) error {
	if goodsID == 0 {
		return errors.New(fmt.Sprintf("产品ID不能为空"))
	}
	if strings.EqualFold(groupName, "") {
		return nil
	}
	hasAttr := service.GetGroupByGoodsIDAndName(db.Orm(), goodsID, groupName) //repository.GoodsAttributesGroup.GetByGoodsIDAndName(goodsID, groupName)

	if hasAttr.IsZero() == false {
		return errors.New(fmt.Sprintf("属性名：%v已经存在", groupName))
	}
	err := dao.Create(db.Orm(), &model.GoodsAttributesGroup{
		GoodsID: goodsID,
		Name:    groupName,
	})
	return err
}
