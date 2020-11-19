package activity

import (
	"strings"

	"github.com/nbvghost/glog"

	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/gweb/tool"
	"github.com/pkg/errors"
)

type CollageService struct {
	dao.BaseDao
}

func (service CollageService) GetCollageGoodsByGoodsID(GoodsID uint64, OID uint64) dao.CollageGoods {
	var timesellGoods dao.CollageGoods
	dao.Orm().Model(&dao.CollageGoods{}).Where("GoodsID=? and OID=?", GoodsID, OID).First(&timesellGoods)
	return timesellGoods
}

func (service CollageService) DeleteCollage(TimeSellID uint64) error {
	//timesell := TimeSellService{}.GetTimeSellByGoodsID(GoodsID)
	var ts dao.Collage
	service.Get(dao.Orm(), TimeSellID, &ts)
	//err := service.Delete(dao.Orm(), &dao.TimeSell{}, ts.ID)
	err := service.DeleteWhere(dao.Orm(), &dao.Collage{}, "Hash=?", ts.Hash)
	glog.Error(err)

	return err
}

func (service CollageService) GetItemByHash(Hash string) dao.Collage {
	var timesell dao.Collage
	err := dao.Orm().Model(&dao.Collage{}).Where("Hash=?", Hash).First(&timesell).Error
	glog.Error(err)
	return timesell
}

func (service CollageService) GetCollageByGoodsID(GoodsID uint64, OID uint64) dao.Collage {
	//timesellGoods := service.GetTimeSellGoodsByGoodsID(GoodsID, OID)
	var timesellGoods dao.CollageGoods
	dao.Orm().Model(&dao.CollageGoods{}).Where("GoodsID=? and OID=?", GoodsID, OID).First(&timesellGoods)

	var timesell dao.Collage
	dao.Orm().Model(&dao.Collage{}).Where("Hash=? and OID=?", timesellGoods.CollageHash, timesellGoods.OID).First(&timesell)
	return timesell

	/*var timesell dao.Collage
	err := dao.Orm().Model(&dao.Collage{}).Where("GoodsID=?", GoodsID).First(&timesell).Error
	glog.Error(err)
	return timesell*/
}
func (service CollageService) GetCollageByHash(Hash string, OID uint64) dao.Collage {
	var timesell dao.Collage
	err := dao.Orm().Model(&dao.Collage{}).Where("Hash=? and OID=?", Hash, OID).First(&timesell).Error
	glog.Error(err)
	return timesell
}
func (service CollageService) AddCollageRecord(OrderNo, OrdersGoodsNo, No string, UserID uint64) error {
	cr := &dao.CollageRecord{}
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
	return service.Add(dao.Orm(), cr)
}
func (service CollageService) FindCollageRecordByUserIDAndNo(UserID uint64, No string) dao.CollageRecord {
	Orm := dao.Orm()
	var cr dao.CollageRecord
	Orm.Model(&dao.CollageRecord{}).Where("UserID=? and No=?").First(&cr)
	return cr

}
