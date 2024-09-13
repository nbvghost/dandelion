package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service/internal/activity"
	"github.com/nbvghost/dandelion/service/internal/company"
	"github.com/nbvghost/dandelion/service/internal/goods"
	"github.com/nbvghost/dandelion/service/internal/journal"
	"github.com/nbvghost/dandelion/service/internal/user"
	"log"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/play"
	"github.com/nbvghost/dandelion/library/result"
	"gorm.io/gorm"
)

type VerificationService struct {
	model.BaseDao
	CardItem   activity.CardItemService
	Journal    journal.JournalService
	Settlement activity.SettlementService
	Orders     OrdersService
	ScoreGoods activity.ScoreGoodsService
	Voucher    activity.VoucherService
	User       user.UserService
	Store      company.StoreService
	Goods      goods.GoodsService
}

// 核销卡卷
func (service VerificationService) VerificationCardItem(DB *gorm.DB, VerificationNo string, Quantity uint, user *model.User, store *model.Store) error {

	verification := service.GetVerificationByVerificationNo(VerificationNo)
	if verification.ID == 0 {
		return errors.New("找不到核销单")
	}

	//var cardItem model.CardItem
	cardItem := dao.GetByPrimaryKey(DB, entity.CardItem, verification.CardItemID).(*model.CardItem)
	if cardItem.IsZero() {
		return gorm.ErrRecordNotFound
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

	err := dao.UpdateByPrimaryKey(DB, entity.Verification, verification.ID, &model.Verification{StoreID: store.ID, StoreUserID: user.ID, Quantity: Quantity})
	if err != nil {
		return err
	}
	err = dao.UpdateByPrimaryKey(DB, entity.CardItem, cardItem.ID, &model.CardItem{UseQuantity: cardItem.UseQuantity + Quantity})
	if err != nil {
		return err
	}

	if strings.EqualFold(cardItem.Type, play.CardItem_Type_OrdersGoods) {
		//如果是商品订单结算门店佣金/结算用户上下级佣金
		//var ordersGoods model.OrdersGoods
		ordersGoods := dao.GetByPrimaryKey(DB, entity.OrdersGoods, cardItem.OrdersGoodsID).(*model.OrdersGoods)
		//err := util.StringToJSON(cardItem.Data, &ordersGoods)
		if ordersGoods.IsZero() {
			return gorm.ErrRecordNotFound
		}
		//var orders model.Orders

		orders := dao.GetByPrimaryKey(DB, entity.Orders, ordersGoods.OrdersID).(*model.Orders)
		if orders.ID == 0 {
			return errors.New("找不到订单或无效订单")
		}
		//CardItem

		var goods model.Goods = ordersGoods.Goods
		/*err = util.JSONToStruct(ordersGoods.Goods, &goods)
		if err != nil {
			return err
		}*/

		go func() {
			//var _goods model.Goods
			_goods := dao.GetByPrimaryKey(db.Orm(), entity.Goods, goods.ID).(*model.Goods)
			if _goods.ID != 0 {
				dao.UpdateByPrimaryKey(db.Orm(), entity.Goods, _goods.ID, &model.Goods{CountSale: _goods.CountSale + uint(Quantity)})
			}
		}()

		var specification model.Specification = ordersGoods.Specification
		/*err = util.JSONToStruct(ordersGoods.Specification, &specification)
		if err != nil {
			return err
		}*/

		if orders.PostType == 1 {
			//邮寄订单，给利润给
			_, err = service.Journal.AddStoreJournal(DB, store.ID, "商品核销", goods.Title+"("+specification.Label+")", play.StoreJournal_Type_HX, int64(specification.MarketPrice-specification.CostPrice), cardItem.ID)
			if err != nil {
				return err
			}
		} else if orders.PostType == 2 {
			//线下订单，给成本价
			_, err = service.Journal.AddStoreJournal(DB, store.ID, "商品核销", goods.Title+"("+specification.Label+")", play.StoreJournal_Type_HX, int64(specification.MarketPrice-specification.CostPrice), cardItem.ID)
			if err != nil {
				return err
			}
			//要减掉门店的库存
			storeStock := service.Store.GetByGoodsIDAndSpecificationIDAndStoreID(goods.ID, specification.ID, store.ID)
			if int64(storeStock.Stock-storeStock.UseStock-uint(Quantity)) < 0 {
				return errors.New("门店库存不足，无法核销")
			} else {
				err = dao.UpdateByPrimaryKey(DB, entity.StoreStock, storeStock.ID, &model.StoreStock{UseStock: storeStock.UseStock + uint(Quantity)})
				if err != nil {
					return err
				}
			}

			if !(orders.Status == model.OrdersStatusOrderOk) {
				//当线下订单被核销后，主订单完成，并产生用户的结算
				err = dao.UpdateByPrimaryKey(DB, entity.Orders, orders.ID, &model.Orders{Status: model.OrdersStatusOrderOk, ReceiptTime: time.Now()})
				if err != nil {
					return err
				}

				ogs, err := service.Orders.FindOrdersGoodsByOrdersID(DB, orders.ID)
				if err != nil {
					return err
				}

				var ogsList []*model.OrdersGoods
				for i := range ogs {
					value := ogs[i]
					//var specification model.Specification
					//util.JSONToStruct(value.Specification, &specification)
					ogsList = append(ogsList, value)
				}

				//线下订单，由核销后结算给用户。邮寄快递，由确定收货时，结算。
				err = service.Settlement.SettlementUser(DB, ogsList, orders)
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
			var scoreGoods model.ScoreGoods
			//err := util.JSONToMap(cardItem.Data, &scoreGoods)
			err = json.Unmarshal([]byte(cardItem.Data),&scoreGoods)
			if err != nil {
				return err
			}
			//err := service.ScoreGoods.Get(DB, cardItem.ScoreGoodsID, &scoreGoods)

			_, err = service.Journal.AddStoreJournal(DB, store.ID, "积分商品核销", scoreGoods.Name, play.StoreJournal_Type_SG, int64(scoreGoods.Price), cardItem.ID)
			if err != nil {
				return err
			}
		} else if strings.EqualFold(cardItem.Type, play.CardItem_Type_Voucher) {
			var voucher model.Voucher
			//err := service.Voucher.Get(DB, cardItem.VoucherID, &voucher)
			//err := util.JSONToStruct(cardItem.Data, &voucher)
			err = json.Unmarshal([]byte(cardItem.Data),&voucher)
			if err != nil {
				return err
			}
			_, err = service.Journal.AddStoreJournal(DB, store.ID, "福利卷核销", voucher.Name, play.StoreJournal_Type_FL, int64(voucher.Amount), cardItem.ID)
			if err != nil {
				return err
			}
		} else {
			return errors.New("未知卡卷类型")
		}

	}

	return nil
}
func (service VerificationService) GetVerificationByVerificationNo(VerificationNo string) model.Verification {
	Orm := db.Orm()
	item := model.Verification{}
	err := Orm.Where("VerificationNo=?", VerificationNo).First(&item).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	log.Println(err)
	return item
}

func (service VerificationService) VerificationSelf(StoreID, StoreStockID dao.PrimaryKey, Quantity uint) *result.ActionResult {
	Orm := db.Orm()
	//var ss model.StoreStock
	ss := dao.GetByPrimaryKey(Orm, entity.StoreStock, StoreStockID).(*model.StoreStock)
	if ss.IsZero() {
		return (&result.ActionResult{}).SmartError(gorm.ErrRecordNotFound, "", 0)
	}

	if ss.StoreID != StoreID {
		return (&result.ActionResult{}).SmartError(errors.New("找不到相关库存"), "", 0)
	}

	if Quantity == 0 {
		return (&result.ActionResult{}).SmartError(errors.New("无效的数量"), "", 0)
	}

	if Quantity > ss.Stock-ss.UseStock {
		//库存不足
		return (&result.ActionResult{}).SmartError(errors.New("库存不足"), "", 0)
	} else {

		//var specification model.Specification
		specification := dao.GetByPrimaryKey(Orm, entity.Specification, ss.SpecificationID).(*model.Specification)

		//var store model.Store
		store := dao.GetByPrimaryKey(Orm, entity.Store, StoreID).(*model.Store)

		//var goods model.Goods
		goods := dao.GetByPrimaryKey(Orm, entity.Goods, ss.GoodsID).(*model.Goods)

		if ss.GoodsID != specification.GoodsID {
			return (&result.ActionResult{}).SmartError(errors.New("找不到相关库存"), "", 0)
		} else {
			//判断金额,

			if specification.CostPrice*Quantity > store.Amount {
				//金额不足
				return (&result.ActionResult{}).SmartError(errors.New("门店金额不足"), "", specification.CostPrice*Quantity-store.Amount)
			} else {

				tx := Orm.Begin()

				err := dao.UpdateByPrimaryKey(tx, entity.StoreStock, StoreStockID, map[string]interface{}{"UseStock": ss.UseStock + Quantity})
				if err != nil {
					tx.Rollback()
					return (&result.ActionResult{}).SmartError(err, "", 0)
				} else {
					detail := fmt.Sprintf("%v,规格：%v(%v)kg成本价：%v，数量：%v", goods.Title, specification.Label, float64(specification.Num)*float64(specification.Weight)/1000, specification.CostPrice, Quantity)
					_, err = service.Journal.AddStoreJournal(tx, StoreID, "自主核销商品库存", detail, play.StoreJournal_Type_ZZHX, -int64(specification.CostPrice*Quantity), ss.ID)
					if err != nil {
						tx.Rollback()
						return (&result.ActionResult{}).SmartError(err, "", 0)
					} else {
						tx.Commit()
						return (&result.ActionResult{}).SmartError(err, "自主核销成功", 0)
					}

				}

			}
		}

	}

}
