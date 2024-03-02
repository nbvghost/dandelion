package model

import "github.com/nbvghost/dandelion/library/dao"

type CustomizeFieldGroup struct {
	dao.Entity
	OID  dao.PrimaryKey `gorm:"column:OID;index"`
	Name string         `gorm:"column:Name"`
}

func (CustomizeFieldGroup) TableName() string {
	return "CustomizeFieldGroup"
}

type CustomizeField struct {
	dao.Entity
	OID      dao.PrimaryKey `gorm:"column:OID;index"`
	GroupID  dao.PrimaryKey `gorm:"column:GroupID"`
	Type     string         `gorm:"column:Type"` //BLOCK,RICH_NUMERIC_UNIT,FIELD,TEXT,IMAGE
	Field    string         `gorm:"column:Field"`
	Extra    string         `gorm:"column:Extra"` //json Unit  Explain
	ParentID dao.PrimaryKey `gorm:"column:ParentID"`
	Sort     int64          `gorm:"column:Sort"`
	//BlockName string         `gorm:"column:BlockName"`
}

func (CustomizeField) TableName() string {
	return "CustomizeField"
}
