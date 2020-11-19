package journal

import (
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/util"

	"errors"

	"time"

	"github.com/jinzhu/gorm"
)

type JournalService struct {
	dao.BaseDao
	//User user.UserService
}

func (service JournalService) StoreListJournal(StoreID uint64, startDate string, endDate string) []dao.StoreJournal {
	Orm := dao.Orm()
	//startDate
	//endDate

	sd, _ := time.ParseInLocation("2006-1-2", startDate, time.Local)
	ed, _ := time.ParseInLocation("2006-1-2", endDate, time.Local)

	ed = ed.Add(24 * time.Hour)

	var StoreJournals []dao.StoreJournal
	Orm.Where("StoreID=?", StoreID).Where("CreatedAt>=?", sd).Where("CreatedAt<?", ed).Order("CreatedAt desc").Find(&StoreJournals)

	return StoreJournals

}
func (service JournalService) AddStoreJournal(DB *gorm.DB, StoreID uint64, Name, Detail string, Type int, Amount int64, TargetID uint64) error {

	logger := &dao.StoreJournal{}
	logger.Name = Name
	logger.Detail = Detail
	logger.Type = Type
	logger.Amount = Amount
	logger.StoreID = StoreID
	logger.TargetID = TargetID

	var store dao.Store
	DB.First(&store, StoreID)

	Balance := int64(store.Amount) + Amount
	if Balance < 0 {
		return errors.New("余额不足")
	}
	logger.Balance = uint64(Balance)

	err := service.ChangeMap(DB, StoreID, &dao.Store{}, map[string]interface{}{"Amount": Balance})
	if err != nil {
		return err
	}
	return service.Add(DB, logger)
}

func (service JournalService) ListUserJournalLeveBrokerage(UserID uint64, IDs []uint64) interface{} {

	Orm := dao.Orm()
	type Result struct {
		TotalAmount int64 `gorm:"column:TotalAmount"`
		*dao.User   `json:"User"`
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
	//var recordsTotal uint64
	db := Orm.Table("UserJournal").Select("SUM(UserJournal.Amount) as TotalAmount,User.*").Joins("JOIN User on User.ID = UserJournal.FromUserID").Where("UserJournal.FromUserID in (?)", IDs).Where("UserJournal.UserID = ?", UserID).Group("UserJournal.FromUserID")
	//db.Limit(10).Offset(0).Find(&result)
	db.Find(&result)
	//db.Offset(0).Count(&recordsTotal)
	//fmt.Println(result)
	//fmt.Println(recordsTotal)

	return result
}

//OrganizationJournal
func (service JournalService) AddOrganizationJournal(DB *gorm.DB, OID uint64, Name, Detail string, Type int, Amount int64, KV dao.KV) error {

	logger := &dao.OrganizationJournal{}
	logger.Name = Name
	logger.Detail = Detail
	logger.Type = Type
	logger.Amount = Amount
	logger.OID = OID
	logger.DataKV = util.StructToJSON(KV)

	var org dao.Organization
	DB.First(&org, OID)

	Balance := int64(org.Amount) + Amount
	if Balance < 0 {
		return errors.New("余额不足")
	}
	logger.Balance = uint64(Balance)

	err := service.ChangeMap(DB, OID, &dao.Organization{}, map[string]interface{}{"Amount": Balance})
	if err != nil {
		return err
	}
	err = service.Add(DB, logger)

	return err
}

func (service JournalService) AddUserJournal(DB *gorm.DB, UserID uint64, Name, Detail string, Type int, Amount int64, KV dao.KV, FromUserID uint64) error {

	logger := &dao.UserJournal{}
	logger.Name = Name
	logger.Detail = Detail
	logger.Type = Type
	logger.Amount = Amount
	logger.UserID = UserID
	logger.DataKV = util.StructToJSON(KV)
	logger.FromUserID = FromUserID

	var user dao.User
	DB.First(&user, UserID)

	Balance := int64(user.Amount) + Amount
	if Balance < 0 {
		return errors.New("余额不足")
	}
	logger.Balance = uint64(Balance)

	err := service.ChangeMap(DB, UserID, &dao.User{}, map[string]interface{}{"Amount": Balance})
	if err != nil {
		return err
	}
	err = service.Add(DB, logger)

	return err
}
func (service JournalService) AddScoreJournal(DB *gorm.DB, UserID uint64, Name, Detail string, Type int, Score int64, KV dao.KV) error {

	logger := &dao.ScoreJournal{}
	logger.Name = Name
	logger.Detail = Detail
	logger.Type = Type
	logger.Score = Score
	logger.UserID = UserID
	logger.DataKV = util.StructToJSON(KV)

	var user dao.User
	DB.First(&user, UserID)

	Balance := int64(user.Score) + Score
	if Balance < 0 {
		return errors.New("积分不足")
	}
	logger.Balance = uint64(Balance)

	err := service.ChangeMap(DB, UserID, &dao.User{}, map[string]interface{}{"Score": Balance})
	if err != nil {
		return err
	}
	err = service.Add(DB, logger)

	return err
}
