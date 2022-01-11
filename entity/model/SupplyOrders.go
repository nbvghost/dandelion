package model

import (
	"errors"
	"runtime/debug"
	"time"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/base"
	"github.com/nbvghost/gpa/types"
)

//充值
type SupplyOrders struct {
	base.BaseModel
	OID      types.PrimaryKey `gorm:"column:OID"`
	UserID   types.PrimaryKey `gorm:"column:UserID"`         //用户ID，支付的用户ID
	OrderNo  string           `gorm:"column:OrderNo;unique"` //订单号
	StoreID  types.PrimaryKey `gorm:"column:StoreID"`        //目标ID，如果门店充值的话，这个就是门店的ID，如果普通用户充值的话，这个就是用户ID
	Type     string           `gorm:"column:Type"`           //Store/User
	PayMoney uint             `gorm:"column:PayMoney"`       //支付金额
	IsPay    uint             `gorm:"column:IsPay"`          //是否支付成功,0=未支付，1，支付成功，2过期
	PayTime  time.Time        `gorm:"column:PayTime"`        //支付时间
}

func (u *SupplyOrders) BeforeCreate(scope *gorm.DB) (err error) {
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
func (SupplyOrders) TableName() string {
	return "SupplyOrders"
}
