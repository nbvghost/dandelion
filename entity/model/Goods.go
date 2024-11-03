package model

import (
	"encoding/json"
	"errors"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"log"
	"runtime/debug"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/library/dao"
)

type GoodsAttribute struct {
	Name  string
	Value string
}

// 商品
type Goods struct {
	dao.Entity
	OID               dao.PrimaryKey        `gorm:"column:OID;index"`                     //
	Uri               string                `gorm:"column:Uri"`                           //
	Title             string                `gorm:"column:Title"`                         //
	GoodsTypeID       dao.PrimaryKey        `gorm:"column:GoodsTypeID"`                   //
	GoodsTypeChildID  dao.PrimaryKey        `gorm:"column:GoodsTypeChildID"`              //
	Price             uint                  `gorm:"column:Price"`                         //
	Stock             uint                  `gorm:"column:Stock"`                         //
	Hide              uint                  `gorm:"column:Hide"`                          //
	Images            sqltype.Array[string] `gorm:"column:Images;type:JSON;"`             //json array//焦点图片
	Videos            sqltype.Array[string] `gorm:"column:Videos;type:JSON;"`             //json array//介绍图片
	Summary           string                `gorm:"column:Summary;type:text"`             //
	Introduce         string                `gorm:"column:Introduce;type:text"`           //
	Pictures          sqltype.Array[string] `gorm:"column:Pictures;type:JSON"`            //json array
	Params            string                `gorm:"column:Params;type:JSON;default:'[]'"` //json array
	ExpressTemplateID dao.PrimaryKey        `gorm:"column:ExpressTemplateID"`             //
	CountSale         uint                  `gorm:"column:CountSale"`                     //销售量
	CountView         uint                  `gorm:"column:CountView"`                     //查看数量
	OrderMinNum       int                   `gorm:"column:OrderMinNum"`                   //最小订购数量
	Tags              pq.StringArray        `gorm:"column:Tags;type:text[]"`              //
	IsRichText        bool                  `gorm:"column:IsRichText;type:boolean"`       //指明Introduce字段是否使用rich text编辑
	Source            string                `gorm:"column:Source"`                        //标记，数据来源
	//TimeSellID        uint `gorm:"column:TimeSellID"`                          //
}

func (u *Goods) GetParams() any {
	var v any
	err := json.Unmarshal([]byte(u.Params), &v)
	if err != nil {
		log.Println(err)
	}
	return v
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
func (u *Goods) TableName() string {
	return "Goods"
}

type GoodsType struct {
	dao.Entity
	OID          dao.PrimaryKey `gorm:"column:OID;index"`
	Uri          string         `gorm:"column:Uri"`
	Name         string         `gorm:"column:Name"`
	Introduction string         `gorm:"column:Introduction"` //主类介绍
	IsStickyTop  bool           `gorm:"column:IsStickyTop"`  //
	Badge        string         `gorm:"column:Badge"`        //徽章
	Image        string         `gorm:"column:Image"`
	ShowAtMenu   bool           `gorm:"column:ShowAtMenu"`
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
	dao.Entity
	OID         dao.PrimaryKey `gorm:"column:OID;index"`
	Uri         string         `gorm:"column:Uri"`
	Name        string         `gorm:"column:Name"`
	Image       string         `gorm:"column:Image"`
	GoodsTypeID dao.PrimaryKey `gorm:"column:GoodsTypeID"`
}

/*
	func (u *GoodsTypeChild) BeforeCreate(scope *gorm.Scope) (err error) {
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
	}
*/
func (GoodsTypeChild) TableName() string {
	return "GoodsTypeChild"
}
