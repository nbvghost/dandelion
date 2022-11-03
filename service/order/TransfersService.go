package order

import (
	"errors"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/play"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/journal"
	"github.com/nbvghost/dandelion/service/wechat"

	"github.com/nbvghost/gpa/types"
	"github.com/nbvghost/tool"
)

type TransfersService struct {
	model.BaseDao
	Journal journal.JournalService
	Wx      wechat.WxService
}

func (service TransfersService) UserTransfers(UserID types.PrimaryKey, ReUserName, IP string, wxConfig *model.WechatConfig) error {
	Orm := singleton.Orm().Begin()

	//var user model.User
	user := dao.GetByPrimaryKey(Orm, entity.User, UserID).(*model.User)
	if user.IsZero() {
		Orm.Rollback()
		return gorm.ErrRecordNotFound
	}
	if user.Amount <= 0 {
		Orm.Rollback()
		return errors.New("余额不足，无法提现")
	}
	if user.Amount < 1000 {
		Orm.Rollback()
		return errors.New("余额不足10元，暂时没法提现，请继续加油")
	}

	transfers := model.Transfers{}

	transfers.Amount = user.Amount

	transfers.Desc = ReUserName + "提现"

	transfers.ReUserName = ReUserName
	transfers.OrderNo = tool.UUID()
	transfers.UserID = user.ID
	transfers.OpenId = user.OpenID
	transfers.IP = IP
	err := dao.Create(Orm, &transfers)
	if err != nil {
		Orm.Rollback()
		return err
	}
	//DB *gorm.DB, UserID uint, Name, Detail string, Type int, Amount int64, TargetID uint,FromUserID uint
	err = service.Journal.AddUserJournal(Orm, user.ID, "提现", ReUserName+"提现", play.UserJournal_Type_TX, -int64(user.Amount), extends.KV{Key: "TransfersOrderNo", Value: transfers.OrderNo}, 0)
	if err != nil {
		Orm.Rollback()
		return err
	}

	Success, Message := service.Wx.Transfers(transfers, wxConfig)
	if Success == false {
		Orm.Rollback()
		return errors.New(Message)
	} else {
		Orm.Commit()
		return nil
	}

}
func (service TransfersService) StoreTransfers(StoreID types.PrimaryKey, UserID types.PrimaryKey, ReUserName, IP string, wxConfig *model.WechatConfig) error {
	Orm := singleton.Orm().Begin()

	//var store model.Store
	store := dao.GetByPrimaryKey(Orm, entity.Store, StoreID).(*model.Store)
	if store.IsZero() {
		Orm.Rollback()
		return gorm.ErrRecordNotFound
	}

	if store.Amount <= 0 {
		Orm.Rollback()
		return errors.New("金额不足，无法提现")
	}

	//var user model.User
	user := dao.GetByPrimaryKey(Orm, entity.User, UserID).(*model.User)
	if user.IsZero() {
		Orm.Rollback()
		return gorm.ErrRecordNotFound
	}
	transfers := model.Transfers{}
	transfers.StoreID = store.ID
	transfers.Amount = store.Amount

	transfers.Desc = ReUserName + "提现"

	transfers.ReUserName = ReUserName
	transfers.OrderNo = tool.UUID()
	transfers.UserID = user.ID
	transfers.OpenId = user.OpenID
	transfers.IP = IP
	err := dao.Create(Orm, &transfers)
	if err != nil {
		Orm.Rollback()
		return err
	}

	err = service.Journal.AddStoreJournal(Orm, store.ID, "提现", ReUserName+"提现", play.StoreJournal_Type_TX, -int64(store.Amount), transfers.ID)
	if err != nil {
		Orm.Rollback()
		return err
	}

	Success, Message := service.Wx.Transfers(transfers, wxConfig)
	if Success == false {
		Orm.Rollback()
		return errors.New(Message)
	} else {
		Orm.Commit()
		return nil
	}

}
