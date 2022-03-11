package repository

import (
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gpa"
	"github.com/nbvghost/gpa/types"
	"reflect"
)

var GoodsAttributes = gpa.Bind(&GoodsAttributesRepository{}, &model.GoodsAttributes{}).(*GoodsAttributesRepository)

type GoodsAttributesRepository struct {
	gpa.IRepository
	FindByGoodsID       func(goodsID types.PrimaryKey) []*model.GoodsAttributes            `gpa:"AutoCrate"`
	GetByGoodsIDAndName func(goodsID types.PrimaryKey, name string) *model.GoodsAttributes `gpa:"AutoCrate"`
	//UpdateByAge    func(age int, update *params.Update) *result.Update                                   `gpa:"AutoCreate"`
	//GetByTel func(tel string) *entity.User `gpa:"AutoCreate"`
}

func (u *GoodsAttributesRepository) Repository() gpa.IRepository {
	return u.IRepository
}
func (u *GoodsAttributesRepository) QueryGoodsAttributesNameInfo() ([]*extends.GoodsAttributesNameInfo, error) {
	rows, err := u.GetDataBase().Query("select * from (select Name,count(Name) as Num from GoodsAttributes group by Name) as m order by m.Num desc", nil)
	if glog.Error(err) {
		return nil, err
	}
	d, err := gpa.Scans(rows, reflect.TypeOf(new(extends.GoodsAttributesNameInfo)), true)
	if glog.Error(err) {
		return nil, err
	}
	//list := gpa.Rows("select * from (select Value,count(Value) as Num from GoodsAttributes where Name=? group by Value) as m order by m.Num desc", []interface{}{name}, &extends.GoodsAttributesValueInfo{})
	return d.([]*extends.GoodsAttributesNameInfo), err

	//list := gpa.Rows("select * from (select Name,count(Name) as Num from GoodsAttributes group by Name) as m order by m.Num desc", nil, &extends.GoodsAttributesNameInfo{})
	//return list.([]*extends.GoodsAttributesNameInfo)
}
func (u *GoodsAttributesRepository) QueryGoodsAttributesValueInfoByName(name string) ([]*extends.GoodsAttributesValueInfo, error) {
	rows, err := u.GetDataBase().Query("select * from (select Value,count(Value) as Num from GoodsAttributes where Name=? group by Value) as m order by m.Num desc", name)
	if glog.Error(err) {
		return nil, err
	}
	d, err := gpa.Scans(rows, reflect.TypeOf(new(extends.GoodsAttributesValueInfo)), true)
	if glog.Error(err) {
		return nil, err
	}
	//list := gpa.Rows("select * from (select Value,count(Value) as Num from GoodsAttributes where Name=? group by Value) as m order by m.Num desc", []interface{}{name}, &extends.GoodsAttributesValueInfo{})
	return d.([]*extends.GoodsAttributesValueInfo), err
}
