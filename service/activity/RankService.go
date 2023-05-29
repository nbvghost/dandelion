package activity

import (
	"errors"
	"github.com/nbvghost/dandelion/library/db"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
)

type RankService struct {
	model.BaseDao
}

func (service RankService) FindDESC() []model.Rank {
	Orm := db.Orm()
	var ranks []model.Rank
	Orm.Model(&model.Rank{}).Order(`"GrowMaxValue" desc`).Find(&ranks)
	return ranks
}
func (service RankService) AddRank(rank model.Rank) error {
	Orm := db.Orm()
	var hasRank model.Rank
	Orm.Model(&model.Rank{}).Where(&model.Rank{GrowMaxValue: rank.GrowMaxValue}).First(&hasRank)
	if hasRank.ID != 0 {
		return errors.New("添加的等级已经存在")
	}
	Orm.Model(&model.Rank{}).Where(&model.Rank{Title: rank.Title}).First(&hasRank)
	if hasRank.ID != 0 {
		return errors.New("添加的等级已经存在")
	}
	err := dao.Create(Orm, &rank)
	return err

}
