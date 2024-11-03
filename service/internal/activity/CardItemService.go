package activity

import (
	"github.com/nbvghost/dandelion/library/db"
	"time"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/play"
	"github.com/nbvghost/dandelion/library/util"

	"gorm.io/gorm"
)

type CardItemService struct {
	model.BaseDao
	Voucher VoucherService
	//Goods   goods.GoodsService
	//Orders  order.OrdersService
}

func (service CardItemService) ListNewCount(UserID dao.PrimaryKey) (TotalRecords int64) {

	Orm := db.Orm()
	var orders []model.CardItem
	db := Orm.Model(model.CardItem{})

	now := time.Now()
	ts := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	te := ts.Add(24 * time.Hour)

	db = db.Where(`"UpdatedAt">=? and "UpdatedAt"<?`, ts, te)
	db = db.Where(`"UserID"=?`, UserID)

	db.Find(&orders).Count(&TotalRecords)
	return

}
func (service CardItemService) CancelOrdersGoodsCardItem(DB *gorm.DB, UserID dao.PrimaryKey, ogs []dao.IEntity) error {

	for i := range ogs {
		value := ogs[i].(*model.OrdersGoods)
		err := dao.DeleteBy(DB, &model.CardItem{}, map[string]interface{}{
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
func (service CardItemService) FindByUserID(UserID dao.PrimaryKey) []model.CardItem {
	var cards []model.CardItem
	//service.FindWhere(singleton.Orm(), &cards, model.CardItem{UserID: UserID})
	service.FindOrderWhere(db.Orm(), `"UseQuantity" asc,"CreatedAt" desc`, &cards, model.CardItem{UserID: UserID})
	return cards
}

// 添加Voucher
func (service CardItemService) AddVoucherCardItem(DB *gorm.DB, OrderNo string, UserID, VoucherID dao.PrimaryKey) error {

	voucher := dao.GetByPrimaryKey(DB, &model.Voucher{}, VoucherID).(*model.Voucher)

	cardItem := model.CardItem{}
	cardItem.OrderNo = OrderNo
	cardItem.PostType = 0
	cardItem.UserID = UserID
	cardItem.Type = CardItem_Type_Voucher
	cardItem.VoucherID = voucher.ID
	cardItem.Data = util.StructToJSON(voucher)
	cardItem.Quantity = 1
	cardItem.ExpireTime = time.Now().Add(time.Hour * 24 * time.Duration(voucher.UseDay))
	err := dao.Create(DB, &cardItem)
	if err != nil {
		return err
	}
	return nil
}

// 添加Voucher
func (service CardItemService) AddScoreGoodsItem(DB *gorm.DB, UserID, ScoreGoodsID dao.PrimaryKey) error {

	//scoreGoodsService := ScoreGoodsService{}

	//var scoreGoods model.ScoreGoods
	//scoreGoodsService.Get(DB, ScoreGoodsID, &scoreGoods)
	scoreGoods := dao.GetByPrimaryKey(DB, &model.ScoreGoods{}, ScoreGoodsID).(*model.ScoreGoods)

	cardItem := model.CardItem{}
	cardItem.PostType = 0
	cardItem.UserID = UserID
	cardItem.Type = play.CardItem_Type_ScoreGoods
	cardItem.VoucherID = 0
	cardItem.ScoreGoodsID = ScoreGoodsID
	cardItem.Data = util.StructToJSON(scoreGoods)
	cardItem.Quantity = 1
	cardItem.ExpireTime = time.Now().Add(time.Hour * 24 * 365)
	err := dao.Create(DB, &cardItem)
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
		err := dao.Create(DB, &cardItem)
		if err != nil {
			return err
		}
	}
	return nil
}
