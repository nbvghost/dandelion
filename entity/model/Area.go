package model

import "github.com/nbvghost/dandelion/library/dao"

type AreaLevel string

const (
	AreaLevelProvince AreaLevel = "PROVINCE"
	AreaLevelCity     AreaLevel = "CITY"
	AreaLevelArea     AreaLevel = "AREA"
	AreaLevelStreet   AreaLevel = "STREET"
	AreaLevelVillage  AreaLevel = "VILLAGE"
)

// Area 中国全国5级行政区划（省、市、县、镇、村）
//
// code,name,level,pcode
//
// level: 省1，市2，县3，镇4，村5
//
// code: 12位，省2位，市2位，县2位，镇3位，村3位
//
// pcode: 直接父级别的code
type Area struct {
	Code         dao.PrimaryKey `gorm:"COMMENT:Code;NOT NULL;column:Code;PRIMARY_KEY"`
	Name         string         `gorm:"column:Name"`
	Level        AreaLevel      `gorm:"column:Level"`
	ProvinceCode dao.PrimaryKey `gorm:"column:ProvinceCode;index"`
	CityCode     dao.PrimaryKey `gorm:"column:CityCode"`
	AreaCode     dao.PrimaryKey `gorm:"column:AreaCode"`
	StreetCode   dao.PrimaryKey `gorm:"column:StreetCode"`
	VillageCode  dao.PrimaryKey `gorm:"column:VillageCode"`
}

func (m *Area) IsZero() bool {
	return m.Code == 0
}

func (m *Area) Primary() dao.PrimaryKey {
	return m.Code
}

func (m *Area) PrimaryName() string {
	return "Code"
}

func (m *Area) TableName() string {
	return "Area"
}
