package model

import (
	"errors"
	"runtime/debug"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/base"
	"github.com/nbvghost/gpa/types"
)

//商品
type Goods struct {
	types.Entity
	OID              types.PrimaryKey `gorm:"column:OID"`
	Title            string           `gorm:"column:Title"`
	GoodsTypeID      types.PrimaryKey `gorm:"column:GoodsTypeID"`
	GoodsTypeChildID types.PrimaryKey `gorm:"column:GoodsTypeChildID"`
	Price            uint             `gorm:"column:Price"`
	Stock            uint             `gorm:"column:Stock"`
	Hide             uint             `gorm:"column:Hide"`
	Images           string           `gorm:"column:Images;type:text;"` //json array
	Videos           string           `gorm:"column:Videos;type:text;"` //json array
	Summary          string           `gorm:"column:Summary;type:text"`
	Introduce        string           `gorm:"column:Introduce;type:text"`
	Pictures         string           `gorm:"column:Pictures;type:text;"` //json array
	Params           string           `gorm:"column:Params;type:text;"`   //json array
	//TimeSellID        uint `gorm:"column:TimeSellID"`                          //
	ExpressTemplateID types.PrimaryKey `gorm:"column:ExpressTemplateID"` //
	CountSale         uint             `gorm:"column:CountSale"`         //销售量
	CountView         uint             `gorm:"column:CountView"`         //查看数量
	Mark              string           `gorm:"column:Mark"`
}

func (u *Goods) BeforeCreate(scope *gorm.DB) (err error) {
	var gt Goods
	scope.Model(u).Where(map[string]interface{}{
		"OID":   u.OID,
		"Title": u.Title,
	}).Find(&gt)
	if gt.ID != 0 {
		err = errors.New("名字重复")
	}
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))

	}
	return
}
func (u Goods) TableName() string {
	return "Goods"
}

type GoodsType struct {
	types.Entity
	OID   types.PrimaryKey `gorm:"column:OID"`
	Name  string           `gorm:"column:Name"`
	Image string           `gorm:"column:Image"`
}

func (GoodsType) TableName() string {
	return "GoodsType"
}

/*func (u *GoodsType) BeforeCreate(scope *gorm.Scope) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))
		return nil
	}
	return nil
}*/
/*func (u *GoodsType) BeforeSave(scope *gorm.Scope) (err error) {
	var gt GoodsType
	scope.DB().Model(u).Where("OID=?", u.OID).Where("Name=?", u.Name).Find(&gt)
	if gt.ID != 0 {
		err = errors.New("名字重复")
	}
	return
}*/

type GoodsTypeChild struct {
	base.BaseModel
	OID         types.PrimaryKey `gorm:"column:OID"`
	Name        string           `gorm:"column:Name"`
	Image       string           `gorm:"column:Image"`
	GoodsTypeID types.PrimaryKey `gorm:"column:GoodsTypeID"`
}

/*func (u *GoodsTypeChild) BeforeCreate(scope *gorm.Scope) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))
		return nil
	}
	return nil
}*/
func (GoodsTypeChild) TableName() string {
	return "GoodsTypeChild"
}
