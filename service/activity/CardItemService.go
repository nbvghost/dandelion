package activity

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/play"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/library/util"

	"github.com/nbvghost/gpa/types"
	"gorm.io/gorm"
	"time"
)

type CardItemService struct {
	model.BaseDao
	Voucher VoucherService
	//Goods   goods.GoodsService
	//Orders  order.OrdersService
}

func (service CardItemService) ListNewCount(UserID types.PrimaryKey) (TotalRecords int64) {

	Orm := singleton.Orm()
	var orders []model.CardItem
	db := Orm.Model(model.CardItem{})

	now := time.Now()
	ts := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	te := ts.Add(24 * time.Hour)

	db = db.Where("UpdatedAt>=? and UpdatedAt<?", ts, te)
	db = db.Where("UserID=?", UserID)

	db.Find(&orders).Count(&TotalRecords)
	return

}
func (service CardItemService) CancelOrdersGoodsCardItem(DB *gorm.DB, UserID types.PrimaryKey, ogs []model.OrdersGoods) error {

	for _, value := range ogs {
		err := service.DeleteWhere(DB, &model.CardItem{}, map[string]interface{}{
			"UserID":        UserID,
			"Type":          CardItem_Type_OrdersGoods,
			"OrdersGoodsID": value.ID,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
func (service CardItemService) FindByUserID(UserID types.PrimaryKey) []model.CardItem {
	var cards []model.CardItem
	//service.FindWhere(singleton.Orm(), &cards, model.CardItem{UserID: UserID})
	service.FindOrderWhere(singleton.Orm(), "UseQuantity asc,CreatedAt desc", &cards, model.CardItem{UserID: UserID})
	return cards
}

//添加Voucher
func (service CardItemService) AddVoucherCardItem(DB *gorm.DB, OrderNo string, UserID, VoucherID types.PrimaryKey) error {
	var voucher model.Voucher
	service.Voucher.Get(DB, VoucherID, &voucher)

	cardItem := model.CardItem{}
	cardItem.OrderNo = OrderNo
	cardItem.PostType = 0
	cardItem.UserID = UserID
	cardItem.Type = CardItem_Type_Voucher
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
func (service CardItemService) AddScoreGoodsItem(DB *gorm.DB, UserID, ScoreGoodsID types.PrimaryKey) error {

	scoreGoodsService := ScoreGoodsService{}

	var scoreGoods model.ScoreGoods
	scoreGoodsService.Get(DB, ScoreGoodsID, &scoreGoods)

	cardItem := model.CardItem{}
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

func (service CardItemService) AddOrdersGoodsCardItem(DB *gorm.DB, orders model.Orders, OrdersGoodse []model.OrdersGoods) error {

	for _, goods := range OrdersGoodse {
		cardItem := model.CardItem{}
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
