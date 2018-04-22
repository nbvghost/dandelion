package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/nbvghost/gweb/tool"
)

type AdminDao = BaseDao

func (AdminDao) DelAdmin(DB *gorm.DB, ID uint64) error {
	err := DB.Delete(Admin{}, "ID=?", ID).Error
	return err
}
func (AdminDao) FindAdmin(Orm *gorm.DB) []Admin {
	list := []Admin{}

	Orm.Find(&list)

	return list
}

func (AdminDao) AddAdmin(DB *gorm.DB, admin *Admin) error {
	return DB.Create(&admin).Error
}
func (AdminDao) GetAdmin(Orm *gorm.DB, ID uint64) *Admin {
	admin := &Admin{}
	err := Orm.Where("ID=?", ID).First(admin).Error //SelectOne(user, "select * from User where Email=?", Email)
	tool.CheckError(err)
	return admin
}
func (AdminDao) GetAdminByEmail(DB *gorm.DB, Email string) *Admin {
	admin := &Admin{}
	err := DB.Where("Email=?", Email).First(admin).Error
	tool.CheckError(err)
	return admin
}
func (AdminDao) GetAdminByTel(DB *gorm.DB, Tel string) *Admin {
	admin := &Admin{}
	err := DB.Where("Tel=?", Tel).First(admin).Error
	tool.CheckError(err)
	return admin
}
