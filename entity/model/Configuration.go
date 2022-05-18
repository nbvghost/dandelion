package model

import (
	"errors"
	"runtime/debug"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/base"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/gpa/types"
)

type Configuration struct {
	base.BaseModel
	OID types.PrimaryKey         `gorm:"column:OID"`
	K   sqltype.ConfigurationKey `gorm:"column:K"`
	V   string                   `gorm:"column:V"`
}

func (u *Configuration) BeforeCreate(scope *gorm.DB) (err error) {
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
func (Configuration) TableName() string {
	return "Configuration"
}
