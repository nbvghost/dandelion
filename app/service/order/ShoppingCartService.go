package order

import (
	"github.com/nbvghost/dandelion/app/service/activity"
	"github.com/nbvghost/dandelion/app/service/company"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/service/goods"
	"github.com/nbvghost/dandelion/app/util"
	"strconv"

	"github.com/nbvghost/glog"
)

type ShoppingCartService struct {
	dao.BaseDao
	Goods        goods.GoodsService
	TimeSell     activity.TimeSellService
	Organization company.OrganizationService
}

func (service ShoppingCartService) GetGSIDs(UserID uint64, GSIDs []uint64) []dao.ShoppingCart {
	Orm := dao.Orm()
	var scs []dao.ShoppingCart
	//Orm := Orm()
	//return Orm.Where("GSID in (?)", IDs).Where(&dao.ShoppingCart{UserID: UserID}).Delete(&dao.ShoppingCart{}).Error
	err := Orm.Where("GSID in (?) and UserID=?", GSIDs, UserID).Find(&scs).Error
	glog.Error(err)
	return scs
}
func (service ShoppingCartService) GetGSID(UserID uint64, GSID uint64) dao.ShoppingCart {
	Orm := dao.Orm()
	var sc dao.ShoppingCart
	err := Orm.Where("GSID=? and  UserID=?", GSID, UserID).First(&sc).Error
	glog.Error(err)
	return sc
}
func (service ShoppingCartService) UpdateByUserIDAndID(UserID uint64, GSID uint64, Quantity uint) error {
	Orm := dao.Orm()
	_sc := service.GetGSID(UserID, GSID)
	//err := service.Get(Orm, ID, &_sc)
	/*err := service.GetGSID(UserID,GSID)
	if err != nil {
		return err
	}*/

	var specification dao.Specification
	util.JSONToStruct(_sc.Specification, &specification)
	/*err = service.Goods.GetSpecification(_sc.SpecificationID, &specification)
	if err != nil {
		return err
	}*/
	if _sc.Quantity > specification.Stock {
		_sc.Quantity = specification.Stock
	}
	if Quantity > specification.Stock {
		_sc.Quantity = specification.Stock
	} else {
		_sc.Quantity = Quantity
	}
	return Orm.Model(&dao.ShoppingCart{}).Where("GSID=?", GSID).Where("UserID=?", UserID).Update(_sc).Error
}

func (service ShoppingCartService) DeleteByUserIDAndGoodsIDAndSpecificationID(UserID, GoodsID, SpecificationID uint64) error {
	Orm := dao.Orm()
	return Orm.Where(&dao.ShoppingCart{UserID: UserID, GSID: strconv.Itoa(int(GoodsID)) + strconv.Itoa(int(SpecificationID))}).Delete(&dao.ShoppingCart{}).Error
}
func (service ShoppingCartService) DeleteListByIDs(UserID uint64, IDs []uint64) error {
	Orm := dao.Orm()
	return Orm.Where("GSID in (?)", IDs).Where(&dao.ShoppingCart{UserID: UserID}).Delete(&dao.ShoppingCart{}).Error
}
func (service ShoppingCartService) FindShoppingCartByUserID(UserID uint64) []dao.ShoppingCart {
	Orm := dao.Orm()
	var list []dao.ShoppingCart
	err := service.FindWhere(Orm, &list, dao.ShoppingCart{UserID: UserID})
	glog.Error(err)
	return list
}
func (service ShoppingCartService) FindShoppingCartListDetails(UserID uint64) (error, []AnalyseOrdersGoods, uint64) {
	//Orm := Orm()
	ordersService := OrdersService{}

	list := service.FindShoppingCartByUserID(UserID)
	//results := make([]map[string]interface{}, 0)
	//var error error

	oredersGoodsList := make([]dao.OrdersGoods, 0)

	for key := range list {

		oredersGoods := ordersService.createOrdersGoods(list[key])
		oredersGoodsList = append(oredersGoodsList, oredersGoods)
		//results[oredersGoods.OID]=append(results[oredersGoods.OID],oredersGoods)
	}

	return ordersService.AnalyseOrdersGoodsList(UserID, dao.Address{}, 0, oredersGoodsList)

}
