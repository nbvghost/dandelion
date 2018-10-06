package service

import (
	"dandelion/app/service/dao"
	"dandelion/app/util"
	"strconv"

	"dandelion/app/play"
	"errors"
	"time"

	"github.com/nbvghost/gweb"
)

type ScoreGoodsService struct {
	dao.BaseDao
	User    UserService
	Journal JournalService
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
func (service ScoreGoodsService) AddScoreGoods(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	item := &dao.ScoreGoods{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	err = service.Add(Orm, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "添加成功", nil)}
}
func (service ScoreGoodsService) GetScoreGoods(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.ScoreGoods{}
	err := service.Get(Orm, ID, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", item)}
}
func (service ScoreGoodsService) DatatablesScoreGoods(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]dao.ScoreGoods{}, company.ID)
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}
func (service ScoreGoodsService) Exchange(UserID, ScoreGoodsID uint64) error {

	cardItemService := CardItemService{}

	tx := dao.Orm().Begin()
	var scoreGoods dao.ScoreGoods
	service.Get(tx, ScoreGoodsID, &scoreGoods)
	if scoreGoods.ID == 0 {
		tx.Rollback()
		return errors.New("暂时无法兑换")
	}
	var user dao.User

	service.User.Get(tx, UserID, &user)
	if user.ID == 0 {
		tx.Rollback()
		return errors.New("用户暂时无法兑换")
	}
	if uint64(scoreGoods.Score) <= user.Score {

		err := cardItemService.AddScoreGoodsItem(tx, UserID, ScoreGoodsID)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = service.Journal.AddScoreJournal(tx, UserID, "积分兑换", scoreGoods.Name, play.ScoreJournal_Type_DH, -int64(scoreGoods.Score), dao.KV{Key: "ScoreGoodsID", Value: scoreGoods.ID})
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
func (service ScoreGoodsService) ListScoreGoods(OID uint64) []dao.ScoreGoods {
	Orm := dao.Orm()
	var list []dao.ScoreGoods
	service.FindAllByOID(Orm, &list, OID)
	return list
}

func (service ScoreGoodsService) DeleteScoreGoods(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.ScoreGoods{}
	err := service.Delete(Orm, item, ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}
func (service ScoreGoodsService) ChangeScoreGoods(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.ScoreGoods{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	err = service.ChangeModel(Orm, ID, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
}
