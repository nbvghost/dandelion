package activity

import (
	"errors"
	"time"

	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/play"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/journal"

	"github.com/nbvghost/gpa/types"
)

type ScoreGoodsService struct {
	model.BaseDao
	//User    user.UserService
	Journal  journal.JournalService
	CardItem CardItemService
}

func (service ScoreGoodsService) Exchange(user *model.User, ScoreGoodsID types.PrimaryKey) error {

	tx := singleton.Orm().Begin()
	var scoreGoods model.ScoreGoods
	service.Get(tx, ScoreGoodsID, &scoreGoods)
	if scoreGoods.ID == 0 {
		tx.Rollback()
		return errors.New("暂时无法兑换")
	}
	if user.ID == 0 {
		tx.Rollback()
		return errors.New("用户暂时无法兑换")
	}
	if uint(scoreGoods.Score) <= user.Score {

		err := service.CardItem.AddScoreGoodsItem(tx, user.ID, ScoreGoodsID)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = service.Journal.AddScoreJournal(tx, user.ID, "积分兑换", scoreGoods.Name, play.ScoreJournal_Type_DH, -int64(scoreGoods.Score), extends.KV{Key: "ScoreGoodsID", Value: scoreGoods.ID})
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

	Orm := singleton.Orm()

	type Result struct {
		TotalScore uint `gorm:"column:TotalScore"`
		TotalCount uint `gorm:"column:TotalCount"`
	}

	//select COUNT(ID),SUM(JSON_EXTRACT(Data,'$.Amount')) as TotalMoney from carditem where Type='Voucher' group by Type;

	var result Result

	Orm.Table("CardItem").Select("SUM(JSON_EXTRACT(Data,'$.Score')) as TotalScore,COUNT(ID) as TotalCount").Where("CreatedAt>=?", st).Where("CreatedAt<?", et).Where("Type=?", "ScoreGoods").Find(&result)
	//fmt.Println(result)
	return result
}

func (service ScoreGoodsService) ListScoreGoods() []model.ScoreGoods {
	Orm := singleton.Orm()
	var list []model.ScoreGoods
	service.FindAll(Orm, &list)
	return list
}
