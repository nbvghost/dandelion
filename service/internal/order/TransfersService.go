package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service/internal/journal"
	"github.com/nbvghost/dandelion/service/internal/wechat"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/transferbatch"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/play"

	"github.com/nbvghost/tool"
)

type TransfersService struct {
	model.BaseDao
	Journal journal.JournalService
	Wx      wechat.WxService
}

func (m TransfersService) UserTransfers(ctx context.Context, UserID dao.PrimaryKey, ReUserName, IP string, wxConfig *model.WechatConfig) error {
	Orm := db.GetDB(ctx).Begin()

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
	transfers.OID = user.OID
	transfers.Amount = user.Amount

	transfers.Desc = ReUserName + "余额提现"

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
	//err = service.Journal.AddUserJournal(Orm, user.ID, "提现", ReUserName+"提现", model.UserJournal_Type_TX, -int64(user.Amount), extends.KV{Key: "TransfersOrderNo", Value: transfers.OrderNo}, user.ID)
	var userJournal *model.UserJournal
	userJournal, err = m.Journal.AddUserJournal(Orm, user.ID, "提现", ReUserName+"提现", -int64(user.Amount), journal.NewDataTypeTransfers(transfers.OrderNo), user.ID)
	if err != nil {
		Orm.Rollback()
		return err
	}

	err = m.Wx.Transfers(transfers, []transferbatch.TransferDetailInput{{
		OutDetailNo:    core.String(fmt.Sprintf("UserJournal%d", userJournal.ID)),
		TransferAmount: core.Int64(int64(transfers.Amount)),
		TransferRemark: core.String("用户余额提现"),
		Openid:         core.String(transfers.OpenId),
		UserName:       core.String(transfers.ReUserName),
	}}, wxConfig)
	if err != nil {
		Orm.Rollback()
		return err
	} else {
		Orm.Commit()
		return nil
	}

}
func (m TransfersService) StoreTransfers(ctx context.Context, StoreID dao.PrimaryKey, UserID dao.PrimaryKey, ReUserName, IP string, wxConfig *model.WechatConfig) error {
	Orm := db.GetDB(ctx).Begin()

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
	transfers.OID = store.OID
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

	var storeJournal *model.StoreJournal
	storeJournal, err = m.Journal.AddStoreJournal(Orm, store.ID, "提现", ReUserName+"提现", play.StoreJournal_Type_TX, -int64(store.Amount), transfers.ID)
	if err != nil {
		Orm.Rollback()
		return err
	}

	err = m.Wx.Transfers(transfers, []transferbatch.TransferDetailInput{{
		OutDetailNo:    core.String(fmt.Sprintf("StoreJournal%d", storeJournal.ID)),
		TransferAmount: core.Int64(int64(transfers.Amount)),
		TransferRemark: core.String("门店余额提现"),
		Openid:         core.String(transfers.OpenId),
		UserName:       core.String(transfers.ReUserName),
	}},
		wxConfig)
	if err != nil {
		Orm.Rollback()
		return err
	} else {
		Orm.Commit()
		return nil
	}

}
