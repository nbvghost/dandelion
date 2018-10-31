package service

import (
	"dandelion/app/play"
	"dandelion/app/service/dao"
	"errors"

	"github.com/nbvghost/gweb/tool"
)

type TransfersService struct {
	dao.BaseDao
	Journal JournalService
	Wx      WxService
}

func (service TransfersService) UserTransfers(UserID uint64, ReUserName, IP string) error {
	Orm := dao.Orm().Begin()

	var user dao.User
	err := service.Get(Orm, UserID, &user)
	if err != nil {
		Orm.Rollback()
		return err
	}
	if user.Amount <= 0 {
		Orm.Rollback()
		return errors.New("余额不足，无法提现")
	}
	if user.Amount < 1000 {
		Orm.Rollback()
		return errors.New("余额不足10元，暂时没法提现，请继续加油")
	}

	transfers := dao.Transfers{}

	transfers.Amount = user.Amount

	transfers.Desc = ReUserName + "提现"

	transfers.ReUserName = ReUserName
	transfers.OrderNo = tool.UUID()
	transfers.UserID = user.ID
	transfers.OpenId = user.OpenID
	transfers.IP = IP
	err = service.Add(Orm, &transfers)
	if err != nil {
		Orm.Rollback()
		return err
	}
	//DB *gorm.DB, UserID uint64, Name, Detail string, Type int, Amount int64, TargetID uint64,FromUserID uint64
	err = service.Journal.AddUserJournal(Orm, user.ID, "提现", ReUserName+"提现", play.UserJournal_Type_TX, -int64(user.Amount), dao.KV{Key: "TransfersOrderNo", Value: transfers.OrderNo}, 0)
	if err != nil {
		Orm.Rollback()
		return err
	}

	Success, Message := service.Wx.Transfers(transfers)
	if Success == false {
		Orm.Rollback()
		return errors.New(Message)
	} else {
		Orm.Commit()
		return nil
	}

}
func (service TransfersService) StoreTransfers(StoreID uint64, UserID uint64, ReUserName, IP string) error {
	Orm := dao.Orm().Begin()

	var store dao.Store
	err := service.Get(Orm, StoreID, &store)
	if err != nil {
		Orm.Rollback()
		return err
	}

	if store.Amount <= 0 {
		Orm.Rollback()
		return errors.New("金额不足，无法提现")
	}

	var user dao.User
	err = service.Get(Orm, UserID, &user)
	if err != nil {
		Orm.Rollback()
		return err
	}
	transfers := dao.Transfers{}
	transfers.StoreID = store.ID
	transfers.Amount = store.Amount

	transfers.Desc = ReUserName + "提现"

	transfers.ReUserName = ReUserName
	transfers.OrderNo = tool.UUID()
	transfers.UserID = user.ID
	transfers.OpenId = user.OpenID
	transfers.IP = IP
	err = service.Add(Orm, &transfers)
	if err != nil {
		Orm.Rollback()
		return err
	}

	err = service.Journal.AddStoreJournal(Orm, store.ID, "提现", ReUserName+"提现", play.StoreJournal_Type_TX, -int64(store.Amount), transfers.ID)
	if err != nil {
		Orm.Rollback()
		return err
	}

	Success, Message := service.Wx.Transfers(transfers)
	if Success == false {
		Orm.Rollback()
		return errors.New(Message)
	} else {
		Orm.Commit()
		return nil
	}

}
