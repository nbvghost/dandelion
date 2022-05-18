package model

import (
	"errors"
	"runtime/debug"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/base"
	"github.com/nbvghost/gpa/types"
)

//门店
type Store struct {
	base.BaseModel
	OID          types.PrimaryKey `gorm:"column:OID"`
	Name         string           `gorm:"column:Name"`
	Address      string           `gorm:"column:Address"`
	Latitude     float64          `gorm:"column:Latitude"`
	Longitude    float64          `gorm:"column:Longitude"`
	Phone        string           `gorm:"column:Phone"`
	Amount       uint             `gorm:"column:Amount"` //现金
	ServicePhone string           `gorm:"column:ServicePhone"`
	OrderPhone   string           `gorm:"column:OrderPhone"`
	ContactName  string           `gorm:"column:ContactName"`
	Introduce    string           `gorm:"column:Introduce"`
	Images       string           `gorm:"column:Images;type:text"`
	Pictures     string           `gorm:"column:Pictures;type:text"`
	Stars        uint             `gorm:"column:Stars"`      //总星星数量
	StarsCount   uint             `gorm:"column:StarsCount"` //评分人数
}

func (u *Store) BeforeCreate(scope *gorm.DB) (err error) {
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
func (Store) TableName() string {
	return "Store"
}
