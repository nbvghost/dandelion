package activity

import (
	"errors"
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/service/journal"
	"time"
)

type ScoreGoodsService struct {
	dao.BaseDao
	//User    user.UserService
	Journal  journal.JournalService
	CardItem CardItemService
}

func (service ScoreGoodsService) Exchange(user *dao.User, ScoreGoodsID uint64) error {

	tx := dao.Orm().Begin()
	var scoreGoods dao.ScoreGoods
	service.Get(tx, ScoreGoodsID, &scoreGoods)
	if scoreGoods.ID == 0 {
		tx.Rollback()
		return errors.New("暂时无法兑换")
	}
	if user.ID == 0 {
		tx.Rollback()
		return errors.New("用户暂时无法兑换")
	}
	if uint64(scoreGoods.Score) <= user.Score {

		err := service.CardItem.AddScoreGoodsItem(tx, user.ID, ScoreGoodsID)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = service.Journal.AddScoreJournal(tx, user.ID, "积分兑换", scoreGoods.Name, play.ScoreJournal_Type_DH, -int64(scoreGoods.Score), dao.KV{Key: "ScoreGoodsID", Value: scoreGoods.ID})
		if err != nil {
			tx.Rollback()
			return err
		} else {
			tx.Commit()
			return nil
		}

	} else {
		tx.Rollback()
		return errors.New("用户积分不足无法兑换")
	}

}
func (service ScoreGoodsService) Situation(StartTime, EndTime int64) interface{} {

	st := time.Unix(StartTime/1000, 0)
	st = time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location())
	et := time.Unix(EndTime/1000, 0).Add(24 * time.Hour)
	et = time.Date(et.Year(), et.Month(), et.Day(), 0, 0, 0, 0, et.Location())

	Orm := dao.Orm()

	type Result struct {
		TotalScore uint64 `gorm:"column:TotalScore"`
		TotalCount uint64 `gorm:"column:TotalCount"`
	}

	//select COUNT(ID),SUM(JSON_EXTRACT(Data,'$.Amount')) as TotalMoney from carditem where Type='Voucher' group by Type;

	var result Result

	Orm.Table("CardItem").Select("SUM(JSON_EXTRACT(Data,'$.Score')) as TotalScore,COUNT(ID) as TotalCount").Where("CreatedAt>=?", st).Where("CreatedAt<?", et).Where("Type=?", "ScoreGoods").Find(&result)
	//fmt.Println(result)
	return result
}

func (service ScoreGoodsService) ListScoreGoods() []dao.ScoreGoods {
	Orm := dao.Orm()
	var list []dao.ScoreGoods
	service.FindAll(Orm, &list)
	return list
}
