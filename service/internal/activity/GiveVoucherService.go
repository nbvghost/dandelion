package activity

import (
	"context"
	"errors"

	"github.com/nbvghost/dandelion/library/db"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
)

type GiveVoucherService struct {
	model.BaseDao
}

func (service GiveVoucherService) FindDESC(ctx context.Context) []model.GiveVoucher {
	Orm := db.GetDB(ctx)
	var ranks []model.GiveVoucher
	Orm.Model(&model.GiveVoucher{}).Order(`"ScoreMaxValue" desc`).Find(&ranks)
	return ranks
}
func (service GiveVoucherService) FindASC(ctx context.Context) []model.GiveVoucher {
	Orm := db.GetDB(ctx)
	var ranks []model.GiveVoucher
	Orm.Model(&model.GiveVoucher{}).Order(`"ScoreMaxValue" asc`).Find(&ranks)
	return ranks
}
func (service GiveVoucherService) AddItem(ctx context.Context, item model.GiveVoucher) error {
	Orm := db.GetDB(ctx)
	var hasRank model.GiveVoucher
	Orm.Model(&model.GiveVoucher{}).Where(&model.GiveVoucher{ScoreMaxValue: item.ScoreMaxValue}).First(&hasRank)
	if hasRank.ID != 0 {
		return errors.New("条件已经存在")
	}
	err := dao.Create(Orm, &item)
	return err

}
func (service GiveVoucherService) SaveItem(ctx context.Context, item model.GiveVoucher) error {
	Orm := db.GetDB(ctx)
	var hasRank model.GiveVoucher
	Orm.Model(&model.GiveVoucher{}).Where(&model.GiveVoucher{ScoreMaxValue: item.ScoreMaxValue}).First(&hasRank)
	if hasRank.ID != 0 && hasRank.ID != item.ID {

		return errors.New("条件已经存在")
	}
	err := dao.Save(Orm, &item)
	return err

}
