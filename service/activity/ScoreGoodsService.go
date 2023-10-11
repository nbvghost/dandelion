package activity

import (
	"errors"
	"github.com/nbvghost/dandelion/library/db"
	"time"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/service/journal"
)

type ScoreGoodsService struct {
	model.BaseDao
	//User    user.UserService
	Journal  journal.JournalService
	CardItem CardItemService
}

func (service ScoreGoodsService) Exchange(user *model.User, ScoreGoodsID dao.PrimaryKey) error {

	tx := db.Orm().Begin()
	//var scoreGoods model.ScoreGoods
	scoreGoods := dao.GetByPrimaryKey(tx, &model.ScoreGoods{}, ScoreGoodsID).(*model.ScoreGoods)
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

		//err = service.Journal.AddScoreJournal(tx, user.ID, "积分兑换", scoreGoods.Name, model.ScoreJournal_Type_DH, -int64(scoreGoods.Score), extends.KV{Key: "ScoreGoodsID", Value: scoreGoods.ID})
		err = service.Journal.AddScoreJournal(tx, user.ID, "积分兑换", scoreGoods.Name, model.ScoreJournal_Type_DH, -int64(scoreGoods.Score))
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

	Orm := db.Orm()

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

func (service ScoreGoodsService) ListScoreGoods() []dao.IEntity {
	Orm := db.Orm()
	//var list []model.ScoreGoods
	//dao.Find(Orm, &model.ScoreGoods{})
	return dao.Find(Orm, &model.ScoreGoods{}).List()
}
