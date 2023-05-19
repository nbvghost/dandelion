package model

import (
	"github.com/nbvghost/gpa/types"
)

type FullTextSearchType string

const (
	FullTextSearchTypeContent  FullTextSearchType = "content"
	FullTextSearchTypeProducts FullTextSearchType = "product"
)

type FullTextSearch struct {
	types.Entity
	OID           types.PrimaryKey   `gorm:"column:OID;index"`                            //
	TID           types.PrimaryKey   `gorm:"column:TID;index;uniqueIndex:idx_unique_id"`  //
	ContentItemID types.PrimaryKey   `gorm:"column:ContentItemID"`                        //
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
