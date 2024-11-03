package activity

import (
	"github.com/nbvghost/dandelion/library/db"
	"log"
	"strings"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/tool"
	"github.com/pkg/errors"
)

type CollageService struct {
	model.BaseDao
}

func (service CollageService) GetCollageGoodsByGoodsID(GoodsID dao.PrimaryKey, OID dao.PrimaryKey) model.CollageGoods {
	var timesellGoods model.CollageGoods
	db.Orm().Model(&model.CollageGoods{}).Where("GoodsID=? and OID=?", GoodsID, OID).First(&timesellGoods)
	return timesellGoods
}

func (service CollageService) DeleteCollage(TimeSellID dao.PrimaryKey) error {
	//timesell := TimeSellService{}.GetTimeSellByGoodsID(GoodsID)
	//var ts model.Collage
	ts := dao.GetByPrimaryKey(db.Orm(), &model.Collage{}, TimeSellID).(*model.Collage)
	//err := service.Delete(singleton.Orm(), &model.TimeSell{}, ts.ID)
	err := dao.DeleteBy(db.Orm(), &model.Collage{}, map[string]interface{}{
		"Hash": ts.Hash,
	})
	log.Println(err)

	return err
}

func (service CollageService) GetItemByHash(Hash string) model.Collage {
	var timesell model.Collage
	err := db.Orm().Model(&model.Collage{}).Where("Hash=?", Hash).First(&timesell).Error
	log.Println(err)
	return timesell
}

func (service CollageService) GetCollageByGoodsID(GoodsID dao.PrimaryKey, OID dao.PrimaryKey) model.Collage {
	//timesellGoods := service.GetTimeSellGoodsByGoodsID(GoodsID, OID)
	var timesellGoods model.CollageGoods
	db.Orm().Model(&model.CollageGoods{}).Where(`"GoodsID"=? and "OID"=?`, GoodsID, OID).First(&timesellGoods)

	var timesell model.Collage
	db.Orm().Model(&model.Collage{}).Where(`"Hash"=? and "OID"=?`, timesellGoods.CollageHash, timesellGoods.OID).First(&timesell)
	return timesell

	/*var timesell model.Collage
	err := singleton.Orm().Model(&model.Collage{}).Where("GoodsID=?", GoodsID).First(&timesell).Error
	log.Println(err)
	return timesell*/
}
func (service CollageService) GetCollageByHash(Hash string, OID dao.PrimaryKey) *model.Collage {
	var timesell model.Collage
	err := db.Orm().Model(&model.Collage{}).Where("Hash=? and OID=?", Hash, OID).First(&timesell).Error
	log.Println(err)
	return &timesell
}
func (service CollageService) AddCollageRecord(OrderNo, OrdersGoodsNo, No string, UserID dao.PrimaryKey) error {
	cr := &model.CollageRecord{}
	cr.No = No
	cr.OrderNo = OrderNo
	cr.UserID = UserID
	cr.OrdersGoodsNo = OrdersGoodsNo
	if strings.EqualFold(No, "") {
		cr.No = tool.UUID()
		cr.Collager = UserID
	} else {
		cr.No = No
		cr.Collager = 0
		_cr := service.FindCollageRecordByUserIDAndNo(UserID, No)
		if _cr.ID != 0 {
			return errors.New("您已经参加了这个活动，看看其它活动吧！")
		}
	}
	return dao.Create(db.Orm(), cr)
}
func (service CollageService) FindCollageRecordByUserIDAndNo(UserID dao.PrimaryKey, No string) model.CollageRecord {
	Orm := db.Orm()
	var cr model.CollageRecord
	Orm.Model(&model.CollageRecord{}).Where("UserID=? and No=?").First(&cr)
	return cr

}
