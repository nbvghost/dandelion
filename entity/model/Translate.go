package model

import (
	"github.com/nbvghost/gpa/types"
)

type Translate struct {
	types.Entity
	Text     string `gorm:"column:Text;uniqueIndex:Translate_Idx_Text_LangType"`     //
	LangType string `gorm:"column:LangType;uniqueIndex:Translate_Idx_Text_LangType"` //
	LangText string `gorm:"column:LangText"`                                         //
}

func (u Translate) TableName() string {
	return "Translate"
}
