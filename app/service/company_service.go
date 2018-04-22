package service

import (
	"dandelion/app/service/dao"
	"strconv"

	"github.com/jinzhu/gorm"
)

type CompanyService struct {
	dao.CompanyDao
}

func (self CompanyService) GetCompany(ID uint64) *dao.Company {
	target := &dao.Company{}
	self.Get(Orm, ID, target)
	return target
}
func (self CompanyService) AddCompany(DB *gorm.DB, shop *dao.Company) error {
	return self.Add(DB, shop)
}
func (self CompanyService) DelCompany(DB *gorm.DB, ID uint64) error {
	return self.Delete(Orm, dao.Company{}, ID)
}
func (self CompanyService) ChangeCompany(ID uint64, Photos, Categories, Province, City, District, Name, Address, Telephone, Special, Opentime, Avgprice, Introduction, Recommend string) bool {
	shop := &dao.Company{}
	shop.ID = ID
	shop.Photos = Photos
	shop.Categories = Categories
	shop.Province = Province
	shop.City = City
	shop.District = District
	shop.Name = Name
	shop.Address = Address
	shop.Telephone = Telephone
	shop.Special = Special
	shop.Opentime = Opentime
	shop.Avgprice, _ = strconv.Atoi(Avgprice)
	shop.Introduction = Introduction
	shop.Recommend = Recommend
	//return Orm.Save(article).Error
	//err := db.Orm.Save(shop).Error
	err := self.ChangeModel(Orm, ID, shop)
	if err != nil {
		return false
	} else {
		return true
	}

}
