package activity

import (
	"errors"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/singleton"
)

type GiveVoucherService struct {
	model.BaseDao
}

func (service GiveVoucherService) FindDESC() []model.GiveVoucher {
	Orm := singleton.Orm()
	var ranks []model.GiveVoucher
	Orm.Model(&model.GiveVoucher{}).Order("ScoreMaxValue desc").Find(&ranks)
	return ranks
}
func (service GiveVoucherService) FindASC() []model.GiveVoucher {
	Orm := singleton.Orm()
	var ranks []model.GiveVoucher
	Orm.Model(&model.GiveVoucher{}).Order("ScoreMaxValue asc").Find(&ranks)
	return ranks
}
func (service GiveVoucherService) AddItem(item model.GiveVoucher) error {
	Orm := singleton.Orm()
	var hasRank model.GiveVoucher
	Orm.Model(&model.GiveVoucher{}).Where(&model.GiveVoucher{ScoreMaxValue: item.ScoreMaxValue}).First(&hasRank)
	if hasRank.ID != 0 {
		return errors.New("条件已经存在")
	}
	err := dao.Create(Orm, &item)
	return err

}
func (service GiveVoucherService) SaveItem(item model.GiveVoucher) error {
	Orm := singleton.Orm()
	var hasRank model.GiveVoucher
	Orm.Model(&model.GiveVoucher{}).Where(&model.GiveVoucher{ScoreMaxValue: item.ScoreMaxValue}).First(&hasRank)
	if hasRank.ID != 0 && hasRank.ID != item.ID {

		return errors.New("条件已经存在")
	}
	err := dao.Save(Orm, &item)
	return err

}
