package service

import (
	"dandelion/app/service/dao"
	"errors"
)

type GiveVoucherService struct {
	dao.BaseDao
}

func (service GiveVoucherService) FindDESC() []dao.GiveVoucher {
	Orm := dao.Orm()
	var ranks []dao.GiveVoucher
	Orm.Model(&dao.GiveVoucher{}).Order("ScoreMaxValue desc").Find(&ranks)
	return ranks
}
func (service GiveVoucherService) FindASC() []dao.GiveVoucher {
	Orm := dao.Orm()
	var ranks []dao.GiveVoucher
	Orm.Model(&dao.GiveVoucher{}).Order("ScoreMaxValue asc").Find(&ranks)
	return ranks
}
func (service GiveVoucherService) AddItem(item dao.GiveVoucher) error {
	Orm := dao.Orm()
	var hasRank dao.GiveVoucher
	Orm.Model(&dao.GiveVoucher{}).Where(&dao.GiveVoucher{ScoreMaxValue: item.ScoreMaxValue}).First(&hasRank)
	if hasRank.ID != 0 {
		return errors.New("条件已经存在")
	}
	err := service.Add(Orm, &item)
	return err

}
func (service GiveVoucherService) SaveItem(item dao.GiveVoucher) error {
	Orm := dao.Orm()
	var hasRank dao.GiveVoucher
	Orm.Model(&dao.GiveVoucher{}).Where(&dao.GiveVoucher{ScoreMaxValue: item.ScoreMaxValue}).First(&hasRank)
	if hasRank.ID != 0 && hasRank.ID != item.ID {

		return errors.New("条件已经存在")
	}
	err := service.Save(Orm, &item)
	return err

}
