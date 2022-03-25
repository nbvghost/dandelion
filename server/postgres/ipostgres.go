package postgres

import (
	"gorm.io/gorm"

	"github.com/nbvghost/gpa"
)

type IPostgres interface {
	Orm() *gorm.DB
	GPA() gpa.IDataBase
}
