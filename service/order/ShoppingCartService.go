package order

import (
	"log"
	"strconv"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/service/activity"
	"github.com/nbvghost/dandelion/service/company"
	"github.com/nbvghost/dandelion/service/goods"

	"github.com/nbvghost/gpa/types"
)

type ShoppingCartService struct {
	model.BaseDao
	Goods        goods.GoodsService
	TimeSell     activity.TimeSellService
	Organization company.OrganizationService
}

func (service ShoppingCartService) GetGSIDs(UserID types.PrimaryKey, GSIDs []string) []model.ShoppingCart {
	Orm := singleton.Orm()
	var scs []model.ShoppingCart
	//Orm := Orm()
	//return Orm.Where("GSID in (?)", IDs).Where(&model.ShoppingCart{UserID: UserID}).Delete(&model.ShoppingCart{}).Error
	err := Orm.Where(`"GSID" in (?) and "UserID"=?`, GSIDs, UserID).Find(&scs).Error
	log.Println(err)
	return scs
}
func (service ShoppingCartService) GetGSID(UserID types.PrimaryKey, GSID string) model.ShoppingCart {
	Orm := singleton.Orm()
	var sc model.ShoppingCart
	err := Orm.Where(`"GSID"=? and  "UserID"=?`, GSID, UserID).First(&sc).Error
	log.Println(err)
	return sc
}
func (service ShoppingCartService) UpdateByUserIDAndID(UserID types.PrimaryKey, GSID string, Quantity uint) error {
	Orm := singleton.Orm()
	_sc := service.GetGSID(UserID, GSID)
	//err := service.Get(Orm, ID, &_sc)
	/*err := service.GetGSID(UserID,GSID)
	if err != nil {
		return err
	}*/

	var specification model.Specification
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
	return Orm.Model(&model.ShoppingCart{}).Where(`"GSID"=?`, GSID).Where(`"UserID"=?`, UserID).Updates(_sc).Error
}

func (service ShoppingCartService) DeleteByUserIDAndGoodsIDAndSpecificationID(db *gorm.DB, UserID, GoodsID, SpecificationID types.PrimaryKey) error {
	return db.Where(&model.ShoppingCart{UserID: UserID, GSID: strconv.Itoa(int(GoodsID)) + strconv.Itoa(int(SpecificationID))}).Delete(&model.ShoppingCart{}).Error
}
func (service ShoppingCartService) DeleteListByIDs(UserID types.PrimaryKey, IDs []string) error {
	Orm := singleton.Orm()
	return Orm.Where(`"GSID" in (?)`, IDs).Where(&model.ShoppingCart{UserID: UserID}).Delete(&model.ShoppingCart{}).Error
}
func (service ShoppingCartService) FindShoppingCartByUserID(UserID types.PrimaryKey) []model.ShoppingCart {
	Orm := singleton.Orm()
	var list []model.ShoppingCart
	err := service.FindWhere(Orm, &list, model.ShoppingCart{UserID: UserID})
	log.Println(err)
	return list
}
func (service ShoppingCartService) FindShoppingCartListDetails(UserID types.PrimaryKey) (error, []AnalyseOrdersGoods, uint) {
	//Orm := Orm()
	ordersService := OrdersService{}

	list := service.FindShoppingCartByUserID(UserID)
	//results := make([]map[string]interface{}, 0)
	//var error error

	oredersGoodsList := make([]model.OrdersGoods, 0)

	for key := range list {

		oredersGoods := ordersService.createOrdersGoods(list[key])
		oredersGoodsList = append(oredersGoodsList, oredersGoods)
		//results[oredersGoods.OID]=append(results[oredersGoods.OID],oredersGoods)
	}

	return ordersService.AnalyseOrdersGoodsList(UserID, extends.Address{}, 0, oredersGoodsList)

}
