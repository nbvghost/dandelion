package model

import (
	"github.com/lib/pq"
	"github.com/nbvghost/dandelion/library/dao"
)

// GoodsReview helpful
type GoodsReview struct {
	dao.Entity
	GoodsID       dao.PrimaryKey `gorm:"column:GoodsID;index:idxGoodsIDUserID"`
	UserID        dao.PrimaryKey `gorm:"column:UserID;index:idxGoodsIDUserID"`
	Title         string         `gorm:"column:Title"`
	Content       string         `gorm:"column:Content"`
	Images        pq.StringArray `gorm:"column:Images;type:text[]"`
	Portrait      string         `gorm:"column:Portrait"`
	NickName      string         `gorm:"column:NickName"`
	Helpful       uint           `gorm:"column:Helpful"`
	IsBuy         bool           `gorm:"column:IsBuy"`
	Specification string         `gorm:"column:Specification"`
	Rating        uint           `gorm:"column:Rating"`
	Like          uint           `gorm:"column:Like"`
}

func (GoodsReview) TableName() string {
	return "GoodsReview"
}
