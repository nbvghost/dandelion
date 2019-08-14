package service

import (
	"dandelion/app/service/dao"

	"errors"

	"time"

	"dandelion/app/play"
	"strings"

	"dandelion/app/util"

	"github.com/jinzhu/gorm"
	"github.com/nbvghost/glog"
)

type VerificationService struct {
	dao.BaseDao
	CardItem   CardItemService
	Journal    JournalService
	Orders     OrdersService
	ScoreGoods ScoreGoodsService
	Voucher    VoucherService
	User       UserService
	StoreStock StoreStockService
	Goods      GoodsService
}

//核销卡卷
func (service VerificationService) VerificationCardItem(DB *gorm.DB, VerificationNo string, Quantity uint, user *dao.User, store *dao.Store) error {

	verification := service.GetVerificationByVerificationNo(VerificationNo)
	if verification.ID == 0 {
		return errors.New("找不到核销单")
	}

	var cardItem dao.CardItem
	err := service.CardItem.Get(DB, verification.CardItemID, &cardItem)
	if err != nil {
		return err
	}

	if Quantity == 0 {
		return errors.New("核销数量不能为0")
	}

	if (cardItem.Quantity - cardItem.UseQuantity) < Quantity {
		return errors.New("核销失败，数量不足")
	}
	if time.Now().Unix() > cardItem.ExpireTime.Unix() {
		return errors.New("卡卷已经过期")
	}

	if verification.StoreID > 0 && verification.StoreUserID > 0 && verification.Quantity > 0 {
		return errors.New("卡卷已经核销")
	}

	err = service.ChangeModel(DB, verification.ID, &dao.Verification{StoreID: store.ID, StoreUserID: user.ID, Quantity: Quantity})
	if err != nil {
		return err
	}
	err = service.ChangeModel(DB, cardItem.ID, &dao.CardItem{UseQuantity: cardItem.UseQuantity + Quantity})
	if err != nil {
		return err
	}

	if strings.EqualFold(cardItem.Type, play.CardItem_Type_OrdersGoods) {
		//如果是商品订单结算门店佣金/结算用户上下级佣金
		var ordersGoods dao.OrdersGoods
		err = service.Goods.Get(DB, cardItem.OrdersGoodsID, &ordersGoods)
		//err := util.StringToJSON(cardItem.Data, &ordersGoods)
		if err != nil {
			return err
		}
		var orders dao.Orders

		err = service.Orders.Get(DB, ordersGoods.OrdersID, &orders)
		if err != nil {
			return err
		}
		if orders.ID == 0 {
			return errors.New("找不到订单或无效订单")
		}

		//CardItem

		var goods dao.Goods
		err = util.JSONToStruct(ordersGoods.Goods, &goods)
		if err != nil {
			return err
		}

		go func() {
			var _goods dao.Goods
			service.Goods.Get(dao.Orm(), goods.ID, &_goods)
			if _goods.ID != 0 {
				service.Goods.ChangeModel(dao.Orm(), _goods.ID, &dao.Goods{CountSale: _goods.CountSale + uint64(Quantity)})
			}
		}()

		var specification dao.Specification
		err = util.JSONToStruct(ordersGoods.Specification, &specification)
		if err != nil {
			return err
		}

		if orders.PostType == 1 {
			//邮寄订单，给利润给
			err = service.Journal.AddStoreJournal(DB, store.ID, "商品核销", goods.Title+"("+specification.Label+")", play.StoreJournal_Type_HX, int64(specification.MarketPrice-specification.CostPrice), cardItem.ID)
			if err != nil {
				return err
			}
		} else if orders.PostType == 2 {
			//线下订单，给成本价
			err = service.Journal.AddStoreJournal(DB, store.ID, "商品核销", goods.Title+"("+specification.Label+")", play.StoreJournal_Type_HX, int64(specification.MarketPrice-specification.CostPrice), cardItem.ID)
			if err != nil {
				return err
			}
			//要减掉门店的库存
			storeStock := service.StoreStock.GetByGoodsIDAndSpecificationIDAndStoreID(goods.ID, specification.ID, store.ID)
			if int64(storeStock.Stock-storeStock.UseStock-uint64(Quantity)) < 0 {
				return errors.New("门店库存不足，无法核销")
			} else {
				err = service.ChangeModel(DB, storeStock.ID, &dao.StoreStock{UseStock: storeStock.UseStock + uint64(Quantity)})
				if err != nil {
					return err
				}
			}

			if !strings.EqualFold(orders.Status, play.OS_OrderOk) {
				//当线下订单被核销后，主订单完成，并产生用户的结算
				err = service.ChangeModel(DB, orders.ID, &dao.Orders{Status: play.OS_OrderOk, ReceiptTime: time.Now()})
				if err != nil {
					return err
				}

				ogs, err := service.Orders.FindOrdersGoodsByOrdersID(DB, orders.ID)
				if err != nil {
					return err
				}

				var Brokerage uint64
				for _, value := range ogs {
					//var specification dao.Specification
					//util.JSONToStruct(value.Specification, &specification)
					Brokerage = Brokerage + value.TotalBrokerage
				}

				//线下订单，由核销后结算给用户。邮寄快递，由确定收货时，结算。
				err = service.User.SettlementUser(DB, Brokerage, orders)
				if err != nil {
					return err
				}
			}

		} else {
			return errors.New("未知订单配送类型")
		}

	} else {
		//如果是福利卷，结算给门店，没有用户上下级佣金结算
		if strings.EqualFold(cardItem.Type, play.CardItem_Type_ScoreGoods) {
			var scoreGoods dao.ScoreGoods
			err := util.JSONToStruct(cardItem.Data, &scoreGoods)
			//err := service.ScoreGoods.Get(DB, cardItem.ScoreGoodsID, &scoreGoods)
			if err != nil {
				return err
			}

			err = service.Journal.AddStoreJournal(DB, store.ID, "积分商品核销", scoreGoods.Name, play.StoreJournal_Type_SG, int64(scoreGoods.Price), cardItem.ID)
			if err != nil {
				return err
			}
		} else if strings.EqualFold(cardItem.Type, play.CardItem_Type_Voucher) {
			var voucher dao.Voucher
			//err := service.Voucher.Get(DB, cardItem.VoucherID, &voucher)
			err := util.JSONToStruct(cardItem.Data, &voucher)
			if err != nil {
				return err
			}
			err = service.Journal.AddStoreJournal(DB, store.ID, "福利卷核销", voucher.Name, play.StoreJournal_Type_FL, int64(voucher.Amount), cardItem.ID)
			if err != nil {
				return err
			}
		} else {
			return errors.New("未知卡卷类型")
		}

	}

	return nil
}
func (service VerificationService) GetVerificationByVerificationNo(VerificationNo string) dao.Verification {
	Orm := dao.Orm()
	item := dao.Verification{}
	err := Orm.Where("VerificationNo=?", VerificationNo).First(&item).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	glog.Error(err)
	return item
}
