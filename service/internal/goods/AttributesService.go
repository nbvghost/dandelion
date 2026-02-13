package goods

import (
	"context"
	"fmt"
	"strings"

	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
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
func (m AttributesService) GetByGoodsIDAndName(Orm *gorm.DB, oid, goodsID, groupID dao.PrimaryKey, name string) *model.GoodsAttributes {
	item := &model.GoodsAttributes{}
	Orm.Where(`"OID"=? and "GoodsID"=? and "GroupID"=? and "Name"=?`, oid, goodsID, groupID, name).First(item) //SelectOne(user, "select * from User where Tel=?", Tel)
	return item
}

func (m AttributesService) QueryGoodsAttributesNameInfo(ctx context.Context) ([]*extends.GoodsAttributesNameInfo, error) {
	rows, err := db.GetDB(ctx).Raw(`select * from (select "Name",count("Name") as "Num" from "GoodsAttributes" group by "Name") as m order by m."Num" desc`, nil).Rows()
	if err != nil {
		return nil, err
	}
	var list []*extends.GoodsAttributesNameInfo
	err = db.GetDB(ctx).ScanRows(rows, &list)
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
func (m AttributesService) QueryGoodsAttributesValueInfoByName(ctx context.Context, name string) ([]*extends.GoodsAttributesValueInfo, error) {
	rows, err := db.GetDB(ctx).Raw(`select * from (select "Value",count("Value") as "Num" from "GoodsAttributes" where "Name"=? group by "Value") as m order by m."Num" desc`, name).Rows()
	if err != nil {
		return nil, err
	}

	var list []*extends.GoodsAttributesValueInfo
	err = db.GetDB(ctx).ScanRows(rows, &list)
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

func (m AttributesService) AllAttributesName(ctx context.Context) ([]*extends.GoodsAttributesNameInfo, error) {

	return m.QueryGoodsAttributesNameInfo(ctx)
}
func (m AttributesService) AllAttributesByName(ctx context.Context, name string) ([]*extends.GoodsAttributesValueInfo, error) {

	return m.QueryGoodsAttributesValueInfoByName(ctx, name)
}
func (m AttributesService) DeleteGoodsAttributes(ctx context.Context, ID dao.PrimaryKey) error {

	return dao.DeleteByPrimaryKey(db.GetDB(ctx), &model.GoodsAttributes{}, ID) //repository.GoodsAttributes.DeleteByID(ID).Err
}

func (m AttributesService) AddGoodsAttributes(ctx context.Context, oid, goodsID, groupID dao.PrimaryKey, name, value string) (*model.GoodsAttributes, error) {
	if goodsID == 0 || groupID == 0 {
		return nil, errors.New(fmt.Sprintf("产品ID不能为空或组ID不能为空"))
	}
	if strings.EqualFold(name, "") || strings.EqualFold(value, "") {
		return nil, errors.New(fmt.Sprintf("名称不能为空"))
	}
	hasAttr := m.GetByGoodsIDAndName(db.GetDB(ctx), oid, goodsID, groupID, name) //repository.GoodsAttributes.GetByGoodsIDAndName(goodsID, name)
	if hasAttr.IsZero() == false {
		return hasAttr, errors.New(fmt.Sprintf("属性名：%v已经存在", name))
	}
	hasAttr = &model.GoodsAttributes{
		OID:     oid,
		GoodsID: goodsID,
		GroupID: groupID,
		Name:    name,
		Value:   value,
	}
	err := dao.Create(db.GetDB(ctx), hasAttr)
	if err != nil {
		return nil, err
	}
	return hasAttr, nil
}
func (m AttributesService) ListGoodsAttributesGroupByGoodsID(ctx context.Context, goodsID dao.PrimaryKey) []*model.GoodsAttributesGroup {

	return m.FindGroupByGoodsID(db.GetDB(ctx), goodsID) //repository.GoodsAttributesGroup.FindByGoodsID(goodsID)
}
func (m AttributesService) GetGoodsAttributesGroup(ctx context.Context, ID dao.PrimaryKey) dao.IEntity {
	return dao.GetByPrimaryKey(db.GetDB(ctx), &model.GoodsAttributesGroup{}, ID) //repository.GoodsAttributesGroup.GetByID(ID)
}
func (m AttributesService) DeleteGoodsAttributesGroup(ctx context.Context, ID dao.PrimaryKey) error {
	attrs := m.ListGoodsAttributesByGroupID(ctx, ID)

	if len(attrs) > 0 {
		return errors.New(fmt.Sprintf("属性组包含子属性，无法删除"))
	}
	del := dao.DeleteByPrimaryKey(db.GetDB(ctx), &model.GoodsAttributesGroup{}, ID) //repository.GoodsAttributesGroup.DeleteByID(ID)
	return del
}
func (m AttributesService) ListGoodsAttributesByGroupID(ctx context.Context, attributesGroupID dao.PrimaryKey) []*model.GoodsAttributes {
	return m.FindByGroupID(db.GetDB(ctx), attributesGroupID) //repository.GoodsAttributes.FindByGroupID(attributesGroupID)
}
func (m AttributesService) ChangeGoodsAttributesGroup(ctx context.Context, id dao.PrimaryKey, groupName string) error {
	if id == 0 {
		return errors.New(fmt.Sprintf("ID不能为空"))
	}
	if strings.EqualFold(groupName, "") {
		return nil
	}
	hasAttr := m.GetGroupByName(db.GetDB(ctx), groupName) //repository.GoodsAttributesGroup.GetByName(groupName)
	if hasAttr.IsZero() == false {
		return errors.New(fmt.Sprintf("属性名：%v已经存在", groupName))
	}

	err := dao.UpdateByPrimaryKey(db.GetDB(ctx), &model.GoodsAttributesGroup{}, id, map[string]interface{}{"Name": groupName}) //repository.GoodsAttributesGroup.UpdateByID(id, map[string]interface{}{"Name": groupName})

	return err
}
func (m AttributesService) AddGoodsAttributesGroup(ctx context.Context, oid, goodsID dao.PrimaryKey, groupName string) (*model.GoodsAttributesGroup, error) {
	if goodsID == 0 {
		return nil, errors.New(fmt.Sprintf("产品ID不能为空"))
	}
	if strings.EqualFold(groupName, "") {
		return nil, errors.New(fmt.Sprintf("名称不能为空"))
	}
	hasAttr := m.GetGroupByGoodsIDAndName(db.GetDB(ctx), goodsID, groupName) //repository.GoodsAttributesGroup.GetByGoodsIDAndName(goodsID, groupName)

	if hasAttr.IsZero() == false {
		return hasAttr, errors.New(fmt.Sprintf("属性名：%v已经存在", groupName))
	}

	hasAttr = &model.GoodsAttributesGroup{
		OID:     oid,
		GoodsID: goodsID,
		Name:    groupName,
	}
	err := dao.Create(db.GetDB(ctx), hasAttr)
	return hasAttr, err
}
