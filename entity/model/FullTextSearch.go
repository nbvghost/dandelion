package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type FullTextSearchType string

const (
	FullTextSearchTypeContent  FullTextSearchType = "content"
	FullTextSearchTypeProducts FullTextSearchType = "product"
)

type FullTextSearch struct {
	dao.Entity
	OID           dao.PrimaryKey     `gorm:"column:OID;index"`                            //
	TID           dao.PrimaryKey     `gorm:"column:TID;index;uniqueIndex:idx_unique_id"`  //
	ContentItemID dao.PrimaryKey     `gorm:"column:ContentItemID"`                        //
	Title         string             `gorm:"column:Title"`                                //
	Content       string             `gorm:"column:Content;type:text"`                    //
	Picture       string             `gorm:"column:Picture"`                              //
	Type          FullTextSearchType `gorm:"column:Type;index;uniqueIndex:idx_unique_id"` //
	Index         string             `gorm:"column:Index;type:tsvector"`                  //
	Uri           string             `gorm:"column:Uri"`                                  //
}

func (FullTextSearch) TableName() string {
	return "FullTextSearch"
}
