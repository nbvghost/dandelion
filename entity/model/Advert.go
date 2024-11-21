package model

import (
	"github.com/lib/pq"
	"github.com/nbvghost/dandelion/library/dao"
)

type AdvertType string

const (
	AdvertTypePop    AdvertType = "pop"
	AdvertTypeBanner AdvertType = "banner"
	AdvertTypeFloat  AdvertType = "float"
)

type Advert struct {
	dao.Entity
	Matching pq.StringArray `gorm:"column:Matching;type:text[]"` //匹配页面，在哪里页面显示
	Img      string         `gorm:"column:Img"`
	Url      string         `gorm:"column:Url"`
	IsPage   bool           `gorm:"column:IsPage"`
	Show     bool           `gorm:"column:Show"`
	Type     AdvertType     `gorm:"column:Type"`
}

func (Advert) TableName() string {
	return "Advert"
}
