package postgres

import (
	"gorm.io/gorm"
)

type IPostgres interface {
	Orm() *gorm.DB
	//GPA() gpa.IDataBase
}
