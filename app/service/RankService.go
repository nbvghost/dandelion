package service

import (
	"dandelion/app/service/dao"
	"errors"
)

type RankService struct {
	dao.BaseDao
}

func (service RankService) FindDESC() []dao.Rank {
	Orm := dao.Orm()
	var ranks []dao.Rank
	Orm.Model(&dao.Rank{}).Order("GrowMaxValue desc").Find(&ranks)
	return ranks
}
func (service RankService) AddRank(rank dao.Rank) error {
	Orm := dao.Orm()
	var hasRank dao.Rank
	Orm.Model(&dao.Rank{}).Where(&dao.Rank{GrowMaxValue: rank.GrowMaxValue}).First(&hasRank)
	if hasRank.ID != 0 {
		return errors.New("添加的等级已经存在")
	}
	Orm.Model(&dao.Rank{}).Where(&dao.Rank{Title: rank.Title}).First(&hasRank)
	if hasRank.ID != 0 {
		return errors.New("添加的等级已经存在")
	}
	err := service.Add(Orm, &rank)
	return err

}
