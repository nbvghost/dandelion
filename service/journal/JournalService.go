package journal

import (
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/library/db"
	"time"

	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/util"

	"gorm.io/gorm"
)

type JournalService struct {
	model.BaseDao
	//User user.UserService
}

func (service JournalService) StoreListJournal(StoreID dao.PrimaryKey, startDate string, endDate string) []model.StoreJournal {
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
func (service JournalService) AddStoreJournal(DB *gorm.DB, StoreID dao.PrimaryKey, Name, Detail string, Type int, Amount int64, TargetID dao.PrimaryKey) error {

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

type Result struct {
	TotalAmount int64          `gorm:"column:TotalAmount"`
	UserID      dao.PrimaryKey `gorm:"column:UserID"`
	UserName    string         `gorm:"column:UserName"`
}

func (service JournalService) ListUserJournalLeveBrokerage(UserID dao.PrimaryKey, IDs []uint) []Result {
	Orm := db.Orm()

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
	db := Orm.Table(`"UserJournal"`).
		Select(`SUM("UserJournal"."Amount") as "TotalAmount","User"."ID" as "UserID","User"."Name" as "UserName"`).
		Joins(`JOIN "User" on "User"."ID" = "UserJournal"."FromUserID"`).
		Where(`"UserJournal"."FromUserID" in (?)`, IDs).
		Where(`"UserJournal"."UserID" = ?`, UserID).
		Group(`"UserJournal"."FromUserID","User"."ID"`)
	//db.Limit(10).Offset(0).Find(&result)
	db.Find(&result)

	//db.Offset(0).Count(&recordsTotal)
	//fmt.Println(result)
	//fmt.Println(recordsTotal)

	return result
}

// OrganizationJournal
func (service JournalService) AddOrganizationJournal(DB *gorm.DB, OID dao.PrimaryKey, Name, Detail string, Type int, Amount int64, KV extends.KV) error {
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
func (service JournalService) DisableFreezeUserAmount(tx *gorm.DB, UserID dao.PrimaryKey, dataType IDataType, FromUserID dao.PrimaryKey) error {
	//where "DataKV"::json ->> 'Value'='13';
	orm := dao.Find(tx, &model.UserFreezeJournal{}).Where(`"UserID"=?`, UserID).Where(`"FromUserID"=?`, FromUserID).Where(`"FreezeType"=?`, model.FreezeTypeFreeze)
	md := dataType.ToMap()
	for key, value := range md {
		orm.Where(fmt.Sprintf(`"DataKV"::json ->> '%s'='%s'`, key, value))
	}
	list := orm.List()
	if len(list) > 0 {
		item := list[0].(*model.UserFreezeJournal)
		err := dao.UpdateByPrimaryKey(tx, &model.UserFreezeJournal{}, item.ID, map[string]any{"FreezeType": model.FreezeTypeDisable})
		if err != nil {
			return err
		}
	}
	return service.UpdateFreezeUserAmount(tx, UserID)
}
func (service JournalService) UnFreezeUserAmount(tx *gorm.DB, UserID dao.PrimaryKey, dataType IDataType, FromUserID dao.PrimaryKey) error {
	//where "DataKV"::json ->> 'Value'='13';

	orm := dao.Find(tx, &model.UserFreezeJournal{}).Where(`"UserID"=?`, UserID).Where(`"FromUserID"=?`, FromUserID).Where(`"FreezeType"=?`, model.FreezeTypeFreeze)
	md := dataType.ToMap()
	for key, value := range md {
		orm.Where(fmt.Sprintf(`"%s"=?`, key), value)
	}
	list := orm.List()
	if len(list) > 0 {
		item := list[0].(*model.UserFreezeJournal)
		err := dao.UpdateByPrimaryKey(tx, &model.UserFreezeJournal{}, item.ID, map[string]any{"FreezeType": model.FreezeTypeUnFreeze})
		if err != nil {
			return err
		}

		err = service.AddUserJournal(tx, UserID, item.Name, item.Detail, item.Amount, dataType, item.FromUserID)
		if err != nil {
			return err
		}

		//err = service.Journal.AddScoreJournal(Orm, _user.ID, "积分", "佣金积分", model.ScoreJournal_Type_LEVE, int64(leveMenoy), extends.KV{Key: "OrdersID", Value: orders.ID})
		err = service.AddScoreJournal(tx, UserID, "积分", "佣金积分", model.ScoreJournal_Type_LEVE, item.Amount)
		if err != nil {
			return err
		}
	}

	return service.UpdateFreezeUserAmount(tx, UserID)
}
func (service JournalService) UpdateFreezeUserAmount(tx *gorm.DB, UserID dao.PrimaryKey) error {
	var balance int64
	err := dao.Find(tx, &model.UserFreezeJournal{}).Where(`"UserID"=?`, UserID).Where(`"FreezeType"=?`, model.FreezeTypeFreeze).Select(`sum("Amount") as "Amount"`).Scan(&balance)
	if err != nil {
		return err
	}
	err = dao.UpdateByPrimaryKey(tx, &model.User{}, UserID, map[string]interface{}{"BlockAmount": balance})
	if err != nil {
		return err
	}
	return err
}
func (service JournalService) FreezeUserAmount(tx *gorm.DB, UserID dao.PrimaryKey, Name, Detail string, Amount int64, dataType IDataType, FromUserID dao.PrimaryKey) error {
	fromUser := dao.GetByPrimaryKey(tx, &model.User{}, FromUserID).(*model.User)

	logger := &model.UserFreezeJournal{}
	logger.Name = Name
	logger.Detail = Detail
	logger.Type = dataType.GetType()
	logger.Amount = int64(Amount)
	logger.UserID = UserID
	logger.DataKV = dataType.ToJSON()
	logger.FromUserID = fromUser.ID
	logger.FromUserName = fromUser.Name

	err := dao.Create(tx, logger)
	if err != nil {
		return err
	}
	//var user model.User
	//tx.First(&user, UserID)
	//user := dao.GetByPrimaryKey(tx, &model.User{}, UserID).(*model.User)

	return service.UpdateFreezeUserAmount(tx, UserID)

	/*user := dao.GetByPrimaryKey(Orm, &model.User{}, UserID).(*model.User)
	if user.IsZero() {
		return gorm.ErrRecordNotFound
	}

	tm := int64(user.BlockAmount) + Menoy
	if tm < 0 {
		return errors.New("冻结金额不足，无法扣款")
	}

	err := dao.UpdateByPrimaryKey(Orm, &model.User{}, UserID, map[string]interface{}{"BlockAmount": tm})
	return err*/
}

func (service JournalService) AddUserJournal(tx *gorm.DB, UserID dao.PrimaryKey, Name, Detail string, Amount int64, dataType IDataType, FromUserID dao.PrimaryKey) error {
	fromUser := dao.GetByPrimaryKey(tx, &model.User{}, FromUserID).(*model.User)

	logger := &model.UserJournal{}
	logger.Name = Name
	logger.Detail = Detail
	logger.Type = dataType.GetType()
	logger.Amount = Amount
	logger.UserID = UserID
	logger.DataKV = dataType.ToJSON()
	logger.FromUserID = fromUser.ID
	logger.FromUserName = fromUser.Name

	{
		user := dao.GetByPrimaryKey(tx, &model.User{}, UserID).(*model.User)
		balance := int64(user.Amount) + Amount
		if balance < 0 {
			return errors.New("余额不足")
		}
		logger.Balance = uint(balance)
		err := dao.UpdateByPrimaryKey(tx, &model.User{}, UserID, map[string]interface{}{"Amount": balance})
		if err != nil {
			return err
		}
	}
	return dao.Create(tx, logger)
}
func (service JournalService) AddScoreJournal(tx *gorm.DB, UserID dao.PrimaryKey, Name, Detail string, Type model.ScoreJournalType, Score int64) error {
	logger := &model.ScoreJournal{}
	logger.Name = Name
	logger.Detail = Detail
	logger.Type = Type
	logger.Score = Score
	logger.UserID = UserID
	//logger.DataKV = util.StructToJSON(KV)

	//var user model.User
	//DB.First(&user, UserID)
	user := dao.GetByPrimaryKey(tx, &model.User{}, UserID).(*model.User)

	Balance := int64(user.Score) + Score
	if Balance < 0 {
		return errors.New("积分不足")
	}
	logger.Balance = uint(Balance)

	err := dao.UpdateByPrimaryKey(tx, &model.User{}, UserID, map[string]interface{}{"Score": Balance})
	if err != nil {
		return err
	}
	err = dao.Create(tx, logger)
	return err
}
