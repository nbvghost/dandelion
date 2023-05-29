package journal

import (
	"errors"
	"github.com/nbvghost/dandelion/library/db"
	"time"

	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/util"

	"gorm.io/gorm"

	"github.com/nbvghost/gpa/types"
)

type JournalService struct {
	model.BaseDao
	//User user.UserService
}

func (service JournalService) StoreListJournal(StoreID types.PrimaryKey, startDate string, endDate string) []model.StoreJournal {
	Orm := db.Orm()
	//startDate
	//endDate

	sd, _ := time.ParseInLocation("2006-1-2", startDate, time.Local)
	ed, _ := time.ParseInLocation("2006-1-2", endDate, time.Local)

	ed = ed.Add(24 * time.Hour)

	var StoreJournals []model.StoreJournal
	Orm.Where("StoreID=?", StoreID).Where("CreatedAt>=?", sd).Where("CreatedAt<?", ed).Order("CreatedAt desc").Find(&StoreJournals)

	return StoreJournals

}
func (service JournalService) AddStoreJournal(DB *gorm.DB, StoreID types.PrimaryKey, Name, Detail string, Type int, Amount int64, TargetID types.PrimaryKey) error {

	logger := &model.StoreJournal{}
	logger.Name = Name
	logger.Detail = Detail
	logger.Type = Type
	logger.Amount = Amount
	logger.StoreID = StoreID
	logger.TargetID = TargetID

	var store model.Store
	DB.First(&store, StoreID)

	Balance := int64(store.Amount) + Amount
	if Balance < 0 {
		return errors.New("余额不足")
	}
	logger.Balance = uint(Balance)

	err := dao.UpdateByPrimaryKey(DB, &model.Store{}, StoreID, map[string]interface{}{"Amount": Balance})
	if err != nil {
		return err
	}
	return dao.Create(DB, logger)
}

func (service JournalService) ListUserJournalLeveBrokerage(UserID types.PrimaryKey, IDs []uint) interface{} {

	Orm := db.Orm()
	type Result struct {
		TotalAmount int64 `gorm:"column:TotalAmount"`
		*model.User `json:"User"`
	}

	var result []Result

	if len(IDs) == 0 {
		result = make([]Result, 0)
		return result
	}

	/**

	通过UserJournal.FromUserID 获取贡献者ID
	#用户A在线上级或上上级结算佣金的时候，所产生的UserJournal的所有记录 UserJournal.FromUserID 一定是用户A

	#在通过  UserJournal.UserID 来判断获取这个用户的金额，FromUserID与UserID 已经存在的上下级关系。

	#添加了  Where("UserJournal.UserID = ?", UserID)  条件，这样就可以用户相应等级的佣金总金额。

	*/
	//var recordsTotal uint
	db := Orm.Table(`"UserJournal"`).Select(`SUM("UserJournal"."Amount") as "TotalAmount",User.*`).Joins(`JOIN "User" on "User"."ID" = "UserJournal"."FromUserID"`).Where(`"UserJournal"."FromUserID" in (?)`, IDs).Where(`"UserJournal"."UserID" = ?`, UserID).Group(`"UserJournal"."FromUserID"`)
	//db.Limit(10).Offset(0).Find(&result)
	db.Find(&result)
	//db.Offset(0).Count(&recordsTotal)
	//fmt.Println(result)
	//fmt.Println(recordsTotal)

	return result
}

// OrganizationJournal
func (service JournalService) AddOrganizationJournal(DB *gorm.DB, OID types.PrimaryKey, Name, Detail string, Type int, Amount int64, KV extends.KV) error {

	logger := &model.OrganizationJournal{}
	logger.Name = Name
	logger.Detail = Detail
	logger.Type = Type
	logger.Amount = Amount
	logger.OID = OID
	logger.DataKV = util.StructToJSON(KV)

	var org model.Organization
	DB.First(&org, OID)

	Balance := int64(org.Amount) + Amount
	if Balance < 0 {
		return errors.New("余额不足")
	}
	logger.Balance = uint(Balance)

	err := dao.UpdateByPrimaryKey(DB, &model.Organization{}, OID, map[string]interface{}{"Amount": Balance})
	if err != nil {
		return err
	}
	err = dao.Create(DB, logger)

	return err
}

func (service JournalService) AddUserJournal(DB *gorm.DB, UserID types.PrimaryKey, Name, Detail string, Type int, Amount int64, KV extends.KV, FromUserID types.PrimaryKey) error {

	logger := &model.UserJournal{}
	logger.Name = Name
	logger.Detail = Detail
	logger.Type = Type
	logger.Amount = Amount
	logger.UserID = UserID
	logger.DataKV = util.StructToJSON(KV)
	logger.FromUserID = FromUserID

	var user model.User
	DB.First(&user, UserID)

	Balance := int64(user.Amount) + Amount
	if Balance < 0 {
		return errors.New("余额不足")
	}
	logger.Balance = uint(Balance)

	err := dao.UpdateByPrimaryKey(DB, &model.User{}, UserID, map[string]interface{}{"Amount": Balance})
	if err != nil {
		return err
	}
	err = dao.Create(DB, logger)

	return err
}
func (service JournalService) AddScoreJournal(DB *gorm.DB, UserID types.PrimaryKey, Name, Detail string, Type int, Score int64, KV extends.KV) error {

	logger := &model.ScoreJournal{}
	logger.Name = Name
	logger.Detail = Detail
	logger.Type = Type
	logger.Score = Score
	logger.UserID = UserID
	logger.DataKV = util.StructToJSON(KV)

	var user model.User
	DB.First(&user, UserID)

	Balance := int64(user.Score) + Score
	if Balance < 0 {
		return errors.New("积分不足")
	}
	logger.Balance = uint(Balance)

	err := dao.UpdateByPrimaryKey(DB, &model.User{}, UserID, map[string]interface{}{"Score": Balance})
	if err != nil {
		return err
	}
	err = dao.Create(DB, logger)

	return err
}
