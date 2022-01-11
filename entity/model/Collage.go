package model

import (
	"errors"
	"runtime/debug"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/base"
	"github.com/nbvghost/gpa/types"
)

//拼团
type Collage struct {
	base.BaseModel
	OID      types.PrimaryKey `gorm:"column:OID"`
	Hash     string           `gorm:"column:Hash"`     //同一个Hash表示同一个活动
	Num      int              `gorm:"column:Num"`      //拼团人数
	Discount int              `gorm:"column:Discount"` // 折扣
	TotalNum int              `gorm:"column:TotalNum"` //总拼团产品数量
	//GoodsID  uint `gorm:"column:GoodsID"`
}

func (u *Collage) BeforeCreate(scope *gorm.DB) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))

	}
	return nil
}
func (Collage) TableName() string {
	return "Collage"
}
