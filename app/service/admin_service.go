package service

import (
	"dandelion/app/service/dao"
	"dandelion/app/util"
	"time"
)

type AdminService struct {
	dao.AdminDao
}

func (self AdminService) DelAdmin(ID uint64) *dao.ActionStatus {

	tx := Orm.Begin()
	admin := &dao.Admin{}
	self.Get(tx, ID, admin)

	if err := Company.DelCompany(tx, admin.CompanyID); err != nil {
		tx.Rollback()
		return &dao.ActionStatus{false, "删除失败", nil}
	}
	if err := self.Delete(tx, dao.Admin{}, ID); err != nil {
		tx.Rollback()
		return &dao.ActionStatus{false, "删除失败", nil}
	}
	return &dao.ActionStatus{true, "删除成功", nil}
}

func (self AdminService) ChangeAdmin(Name, Password, Email, Tel string, ID uint64) error {

	return self.ChangeModel(Orm, ID, dao.Admin{Name: Name, Password: util.Md5ByString(Password), Email: Email, Tel: Tel})
}

func (self AdminService) AddAdmin(Name, Password, Email, Tel string) *dao.ActionStatus {

	as := &dao.ActionStatus{}

	tx := Orm.Begin()

	admin := &dao.Admin{}
	admin.Name = Name
	admin.Email = Email
	admin.Tel = Tel
	admin.Password = util.Md5ByString(Password)
	admin.LastLoginAt = time.Now()
	if haveAdmin := self.GetAdminByEmail(tx, Email); haveAdmin.ID != 0 {
		tx.Rollback()
		as.Success = false
		as.Message = "这个邮箱已经存在"
		return as
	}

	if haveAdmin := self.GetAdminByTel(tx, Tel); haveAdmin.ID != 0 {
		tx.Rollback()
		as.Success = false
		as.Message = "这个电话已经存在"
		return as
	}

	shop := &dao.Company{}
	shop.Name = Name + "的店铺"
	shop.Expire = time.Now().Add(2 * 24 * time.Hour)
	if err := Company.AddCompany(tx, shop); err != nil {
		tx.Rollback()
		as.Success = false
		as.Message = err.Error()
		return as
	}

	admin.CompanyID = shop.ID

	if err := self.Add(tx, admin); err != nil {
		tx.Rollback()
		as.Success = false
		as.Message = err.Error()
		return as
	}

	tx.Commit()
	as.Success = true
	as.Message = "添加成功"
	return as
}
