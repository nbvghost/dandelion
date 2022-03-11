package activity

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/gpa/types"
	"strings"

	"github.com/nbvghost/glog"

	"github.com/nbvghost/tool"
	"github.com/pkg/errors"
)

type CollageService struct {
	model.BaseDao
}

func (service CollageService) GetCollageGoodsByGoodsID(GoodsID types.PrimaryKey, OID types.PrimaryKey) model.CollageGoods {
	var timesellGoods model.CollageGoods
	singleton.Orm().Model(&model.CollageGoods{}).Where("GoodsID=? and OID=?", GoodsID, OID).First(&timesellGoods)
	return timesellGoods
}

func (service CollageService) DeleteCollage(TimeSellID types.PrimaryKey) error {
	//timesell := TimeSellService{}.GetTimeSellByGoodsID(GoodsID)
	var ts model.Collage
	service.Get(singleton.Orm(), TimeSellID, &ts)
	//err := service.Delete(singleton.Orm(), &model.TimeSell{}, ts.ID)
	err := service.DeleteWhere(singleton.Orm(), &model.Collage{}, map[string]interface{}{
		"Hash": ts.Hash,
	})
	glog.Error(err)

	return err
}

func (service CollageService) GetItemByHash(Hash string) model.Collage {
	var timesell model.Collage
	err := singleton.Orm().Model(&model.Collage{}).Where("Hash=?", Hash).First(&timesell).Error
	glog.Error(err)
	return timesell
}

func (service CollageService) GetCollageByGoodsID(GoodsID types.PrimaryKey, OID types.PrimaryKey) model.Collage {
	//timesellGoods := service.GetTimeSellGoodsByGoodsID(GoodsID, OID)
	var timesellGoods model.CollageGoods
	singleton.Orm().Model(&model.CollageGoods{}).Where("GoodsID=? and OID=?", GoodsID, OID).First(&timesellGoods)

	var timesell model.Collage
	singleton.Orm().Model(&model.Collage{}).Where("Hash=? and OID=?", timesellGoods.CollageHash, timesellGoods.OID).First(&timesell)
	return timesell

	/*var timesell model.Collage
	err := singleton.Orm().Model(&model.Collage{}).Where("GoodsID=?", GoodsID).First(&timesell).Error
	glog.Error(err)
	return timesell*/
}
func (service CollageService) GetCollageByHash(Hash string, OID types.PrimaryKey) model.Collage {
	var timesell model.Collage
	err := singleton.Orm().Model(&model.Collage{}).Where("Hash=? and OID=?", Hash, OID).First(&timesell).Error
	glog.Error(err)
	return timesell
}
func (service CollageService) AddCollageRecord(OrderNo, OrdersGoodsNo, No string, UserID types.PrimaryKey) error {
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
	return service.Add(singleton.Orm(), cr)
}
func (service CollageService) FindCollageRecordByUserIDAndNo(UserID types.PrimaryKey, No string) model.CollageRecord {
	Orm := singleton.Orm()
	var cr model.CollageRecord
	Orm.Model(&model.CollageRecord{}).Where("UserID=? and No=?").First(&cr)
	return cr

}
