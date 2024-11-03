package goods

import (
	"fmt"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
)

type AttributesService struct {
}

////FindByGoodsID       func(goodsID dao.PrimaryKey) ([]*model.GoodsAttributes, error)            `gpa:"AutoCrate"`
//FindByGroupID       func(groupID dao.PrimaryKey) ([]*model.GoodsAttributes, error)            `gpa:"AutoCrate"`
//GetByGoodsIDAndName func(goodsID dao.PrimaryKey, name string) (*model.GoodsAttributes, error) `gpa:"AutoCrate"`

//FindByGoodsID       func(goodsID dao.PrimaryKey) ([]*model.GoodsAttributesGroup, error)            `gpa:"AutoCrate"`
//GetByGoodsIDAndName func(goodsID dao.PrimaryKey, name string) (*model.GoodsAttributesGroup, error) `gpa:"AutoCrate"`
//GetByName           func(name string) (*model.GoodsAttributesGroup, error)                           `gpa:"AutoCrate"`

func (m AttributesService) FindGroupByGoodsID(Orm *gorm.DB, goodsID dao.PrimaryKey) []*model.GoodsAttributesGroup {
	list := make([]*model.GoodsAttributesGroup, 0)
	Orm.Where(`"GoodsID"=?`, goodsID).Find(&list) //SelectOne(user, "select * from User where Tel=?", Tel)
	return list
}
func (m AttributesService) GetGroupByName(Orm *gorm.DB, name string) *model.GoodsAttributesGroup {
	item := &model.GoodsAttributesGroup{}
	Orm.Where(`"Name"=?`, name).First(item) //SelectOne(user, "select * from User where Tel=?", Tel)
	return item
}
func (m AttributesService) GetGroupByGoodsIDAndName(Orm *gorm.DB, goodsID dao.PrimaryKey, name string) *model.GoodsAttributesGroup {
	item := &model.GoodsAttributesGroup{}
	Orm.Where(`"GoodsID"=? and "Name"=?`, goodsID, name).First(item) //SelectOne(user, "select * from User where Tel=?", Tel)
	return item
}

func (m AttributesService) FindByGoodsID(Orm *gorm.DB, goodsID dao.PrimaryKey) []*model.GoodsAttributes {
	list := make([]*model.GoodsAttributes, 0)
	Orm.Where(`"GoodsID"=?`, goodsID).Find(&list) //SelectOne(user, "select * from User where Tel=?", Tel)
	return list
}
func (m AttributesService) FindByGroupID(Orm *gorm.DB, groupID dao.PrimaryKey) []*model.GoodsAttributes {
	list := make([]*model.GoodsAttributes, 0)
	Orm.Where(`"GroupID"=?`, groupID).Find(&list) //SelectOne(user, "select * from User where Tel=?", Tel)
	return list
}
func (m AttributesService) GetByGoodsIDAndName(Orm *gorm.DB, oid, goodsID dao.PrimaryKey, name string) *model.GoodsAttributes {
	item := &model.GoodsAttributes{}
	Orm.Where(`"OID"=? and "GoodsID"=? and "Name"=?`, oid, goodsID, name).First(item) //SelectOne(user, "select * from User where Tel=?", Tel)
	return item
}

func (m AttributesService) QueryGoodsAttributesNameInfo() ([]*extends.GoodsAttributesNameInfo, error) {
	rows, err := db.Orm().Raw(`select * from (select "Name",count("Name") as "Num" from "GoodsAttributes" group by "Name") as m order by m."Num" desc`, nil).Rows()
	if err != nil {
		return nil, err
	}
	var list []*extends.GoodsAttributesNameInfo
	err = db.Orm().ScanRows(rows, &list)
	if err != nil {
		return nil, err
	}
	/*d, err := gpa.ScanRows(rows, reflect.TypeOf(new(extends.GoodsAttributesNameInfo)), true)
	if err != nil {
		return nil, err
	}*/
	//list := gpa.Rows("select * from (select Value,count(Value) as Num from GoodsAttributes where Name=? group by Value) as m order by m.Num desc", []interface{}{name}, &extends.GoodsAttributesValueInfo{})
	return list, err

	//list := gpa.Rows("select * from (select Name,count(Name) as Num from GoodsAttributes group by Name) as m order by m.Num desc", nil, &extends.GoodsAttributesNameInfo{})
	//return list.([]*extends.GoodsAttributesNameInfo)
}
func (m AttributesService) QueryGoodsAttributesValueInfoByName(name string) ([]*extends.GoodsAttributesValueInfo, error) {
	rows, err := db.Orm().Raw(`select * from (select "Value",count("Value") as "Num" from "GoodsAttributes" where "Name"=? group by "Value") as m order by m."Num" desc`, name).Rows()
	if err != nil {
		return nil, err
	}

	var list []*extends.GoodsAttributesValueInfo
	err = db.Orm().ScanRows(rows, &list)
	if err != nil {
		return nil, err
	}

	/*d, err := gpa.ScanRows(rows, reflect.TypeOf(new(extends.GoodsAttributesValueInfo)), true)
	if err != nil {
		return nil, err
	}*/
	//list := gpa.Rows("select * from (select Value,count(Value) as Num from GoodsAttributes where Name=? group by Value) as m order by m.Num desc", []interface{}{name}, &extends.GoodsAttributesValueInfo{})
	return list, err
}

func (m AttributesService) AllAttributesName() ([]*extends.GoodsAttributesNameInfo, error) {

	return m.QueryGoodsAttributesNameInfo()
}
func (m AttributesService) AllAttributesByName(name string) ([]*extends.GoodsAttributesValueInfo, error) {

	return m.QueryGoodsAttributesValueInfoByName(name)
}
func (m AttributesService) DeleteGoodsAttributes(ID dao.PrimaryKey) error {

	return dao.DeleteByPrimaryKey(db.Orm(), &model.GoodsAttributes{}, ID) //repository.GoodsAttributes.DeleteByID(ID).Err
}

func (m AttributesService) AddGoodsAttributes(oid, goodsID, groupID dao.PrimaryKey, name, value string) error {
	if goodsID == 0 || groupID == 0 {
		return errors.New(fmt.Sprintf("产品ID不能为空或组ID不能为空"))
	}
	if strings.EqualFold(name, "") || strings.EqualFold(value, "") {
		return nil
	}
	hasAttr := m.GetByGoodsIDAndName(db.Orm(), oid, goodsID, name) //repository.GoodsAttributes.GetByGoodsIDAndName(goodsID, name)
	if hasAttr.IsZero() == false {
		return errors.New(fmt.Sprintf("属性名：%v已经存在", name))
	}
	err := dao.Create(db.Orm(), &model.GoodsAttributes{
		OID:     oid,
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
func (m AttributesService) ListGoodsAttributesGroupByGoodsID(goodsID dao.PrimaryKey) []*model.GoodsAttributesGroup {

	return m.FindGroupByGoodsID(db.Orm(), goodsID) //repository.GoodsAttributesGroup.FindByGoodsID(goodsID)
}
func (m AttributesService) GetGoodsAttributesGroup(ID dao.PrimaryKey) dao.IEntity {
	return dao.GetByPrimaryKey(db.Orm(), &model.GoodsAttributesGroup{}, ID) //repository.GoodsAttributesGroup.GetByID(ID)
}
func (m AttributesService) DeleteGoodsAttributesGroup(ID dao.PrimaryKey) error {
	attrs := m.ListGoodsAttributesByGroupID(ID)

	if len(attrs) > 0 {
		return errors.New(fmt.Sprintf("属性组包含子属性，无法删除"))
	}
	del := dao.DeleteByPrimaryKey(db.Orm(), &model.GoodsAttributesGroup{}, ID) //repository.GoodsAttributesGroup.DeleteByID(ID)
	return del
}
func (m AttributesService) ListGoodsAttributesByGroupID(attributesGroupID dao.PrimaryKey) []*model.GoodsAttributes {
	return m.FindByGroupID(db.Orm(), attributesGroupID) //repository.GoodsAttributes.FindByGroupID(attributesGroupID)
}
func (m AttributesService) ChangeGoodsAttributesGroup(id dao.PrimaryKey, groupName string) error {
	if id == 0 {
		return errors.New(fmt.Sprintf("ID不能为空"))
	}
	if strings.EqualFold(groupName, "") {
		return nil
	}
	hasAttr := m.GetGroupByName(db.Orm(), groupName) //repository.GoodsAttributesGroup.GetByName(groupName)
	if hasAttr.IsZero() == false {
		return errors.New(fmt.Sprintf("属性名：%v已经存在", groupName))
	}

	err := dao.UpdateByPrimaryKey(db.Orm(), &model.GoodsAttributesGroup{}, id, map[string]interface{}{"Name": groupName}) //repository.GoodsAttributesGroup.UpdateByID(id, map[string]interface{}{"Name": groupName})

	return err
}
func (m AttributesService) AddGoodsAttributesGroup(oid, goodsID dao.PrimaryKey, groupName string) error {
	if goodsID == 0 {
		return errors.New(fmt.Sprintf("产品ID不能为空"))
	}
	if strings.EqualFold(groupName, "") {
		return nil
	}
	hasAttr := m.GetGroupByGoodsIDAndName(db.Orm(), goodsID, groupName) //repository.GoodsAttributesGroup.GetByGoodsIDAndName(goodsID, groupName)

	if hasAttr.IsZero() == false {
		return errors.New(fmt.Sprintf("属性名：%v已经存在", groupName))
	}
	err := dao.Create(db.Orm(), &model.GoodsAttributesGroup{
		OID:     oid,
		GoodsID: goodsID,
		Name:    groupName,
	})
	return err
}
