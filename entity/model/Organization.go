package model

import (
	"time"

	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/gpa/types"
)

type OrganizationStatus int

const (
	OrganizationStatusNormal OrganizationStatus = iota
	OrganizationStatusBlock
)

type Organization struct {
	types.Entity
	//AdminID      types.PrimaryKey    `gorm:"column:AdminID"`
	//Domain       string              `gorm:"column:Domain;not null;unique"`         //三级域名
	Name         string              `gorm:"column:Name;not null"`                  //店名
	Amount       uint                `gorm:"column:Amount;default:0"`               //现金
	BlockAmount  uint                `gorm:"column:BlockAmount;default:0"`          //冻结现金
	Address      string              `gorm:"column:Address"`                        //街道地址
	Telephone    string              `gorm:"column:Telephone"`                      //手机
	Email        string              `gorm:"column:Email"`                          //联系邮箱
	Categories   string              `gorm:"column:Categories"`                     //门店的类型
	Longitude    float64             `gorm:"column:Longitude"`                      //地理位置
	Latitude     float64             `gorm:"column:Latitude"`                       //地理位置
	Photos       sqltype.StringArray `gorm:"column:Photos;type:JSON"`               //店的图片
	Special      sqltype.ObjectArray `gorm:"column:Special;type:JSON;default:'[]'"` //特色
	Opentime     string              `gorm:"column:Opentime"`                       //营业时间
	Avgprice     int                 `gorm:"column:Avgprice"`                       //每人平均消费
	Introduction string              `gorm:"column:Introduction"`                   //介绍
	Recommend    string              `gorm:"column:Recommend"`                      //推荐
	Link         string              `gorm:"column:Link"`                           //链接
	Vip          int                 `gorm:"column:Vip;default:0"`                  //VIP等级
	PayTime      *time.Time          `gorm:"column:PayTime"`                        //缴费时间
	Expire       time.Time           `gorm:"column:Expire"`                         //过期时间
	Status       OrganizationStatus  `gorm:"column:Status"`                         //冻结
	//Province     string    `gorm:"column:Province"`                //省
	//City         string    `gorm:"column:City"`                    //市
	//District     string    `gorm:"column:District"`                //区域
}

func (Organization) TableName() string {
	return "Organization"
}
