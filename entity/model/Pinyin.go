package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type Pinyin struct {
	dao.Entity
	Word   string `gorm:"column:Word;not null;uniqueIndex:Pinyin_Idx_Unique_Word"`
	Pinyin string `gorm:"column:Pinyin"`
	Tone   int    `gorm:"column:Tone;not null"`
}

func (Pinyin) TableName() string {
	return "Pinyin"
}
