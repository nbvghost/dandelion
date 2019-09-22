package service

import (
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/util"
	"time"

	"github.com/jinzhu/gorm"
)

type CardItemService struct {
	dao.BaseDao
	Voucher VoucherService
	Goods   GoodsService
}

func (service CardItemService) ListNewCount(UserID uint64) (TotalRecords int) {

	Orm := dao.Orm()
	var orders []dao.CardItem
	db := Orm.Model(dao.CardItem{})

	now := time.Now()
	ts := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	te := ts.Add(24 * time.Hour)

	db = db.Where("UpdatedAt>=? and UpdatedAt<?", ts, te)
	db = db.Where("UserID=?", UserID)

	db.Find(&orders).Count(&TotalRecords)
	return

}
func (service CardItemService) CancelOrdersGoodsCardItem(DB *gorm.DB, UserID, OrdersID uint64) error {

	ogs, err := GlobalService.Orders.FindOrdersGoodsByOrdersID(DB, OrdersID)

	if err != nil {
		return err
	}

	for _, value := range ogs {
		err := service.DeleteWhere(DB, &dao.CardItem{}, "UserID=? and Type=? and OrdersGoodsID=?", UserID, play.CardItem_Type_OrdersGoods, value.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
func (service CardItemService) FindByUserID(UserID uint64) []dao.CardItem {
	var cards []dao.CardItem
	//service.FindWhere(dao.Orm(), &cards, dao.CardItem{UserID: UserID})
	service.FindOrderWhere(dao.Orm(), "UseQuantity asc,CreatedAt desc", &cards, dao.CardItem{UserID: UserID})
	return cards
}

//添加Voucher
func (service CardItemService) AddVoucherCardItem(DB *gorm.DB, OrderNo string, UserID, VoucherID uint64) error {
	var voucher dao.Voucher
	service.Voucher.Get(DB, VoucherID, &voucher)

	cardItem := dao.CardItem{}
	cardItem.OrderNo = OrderNo
	cardItem.PostType = 0
	cardItem.UserID = UserID
	cardItem.Type = play.CardItem_Type_Voucher
	cardItem.VoucherID = voucher.ID
	cardItem.Data = util.StructToJSON(voucher)
	cardItem.Quantity = 1
	cardItem.ExpireTime = time.Now().Add(time.Hour * 24 * time.Duration(voucher.UseDay))
	err := service.Add(DB, &cardItem)
	if err != nil {
		return err
	}
	return nil
}

//添加Voucher
func (service CardItemService) AddScoreGoodsItem(DB *gorm.DB, UserID, ScoreGoodsID uint64) error {

	scoreGoodsService := ScoreGoodsService{}

	var scoreGoods dao.ScoreGoods
	scoreGoodsService.Get(DB, ScoreGoodsID, &scoreGoods)

	cardItem := dao.CardItem{}
	cardItem.PostType = 0
	cardItem.UserID = UserID
	cardItem.Type = play.CardItem_Type_ScoreGoods
	cardItem.VoucherID = 0
	cardItem.ScoreGoodsID = ScoreGoodsID
	cardItem.Data = util.StructToJSON(scoreGoods)
	cardItem.Quantity = 1
	cardItem.ExpireTime = time.Now().Add(time.Hour * 24 * 365)
	err := service.Add(DB, &cardItem)
	if err != nil {
		return err
	}
	return nil
}

//OrdersGoods,Voucher,ScoreGoods

func (service CardItemService) AddOrdersGoodsCardItem(DB *gorm.DB, orders dao.Orders, OrdersGoodse []dao.OrdersGoods) error {

	for _, goods := range OrdersGoodse {
		cardItem := dao.CardItem{}
		cardItem.PostType = orders.PostType
		cardItem.UserID = orders.UserID
		cardItem.Type = play.CardItem_Type_OrdersGoods
		cardItem.OrdersGoodsID = goods.ID
		cardItem.Data = util.StructToJSON(goods)
		cardItem.Quantity = goods.Quantity
		cardItem.ExpireTime = time.Now().Add(time.Hour * 24 * 365)
		err := service.Add(DB, &cardItem)
		if err != nil {
			return err
		}
	}
	return nil
}
