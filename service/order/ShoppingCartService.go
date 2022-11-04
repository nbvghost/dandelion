package order

import (
	"log"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
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
func (service ShoppingCartService) GetByGoodsIDAndSpecificationID(db *gorm.DB, UserID types.PrimaryKey, GoodsID, SpecificationID types.PrimaryKey) model.ShoppingCart {
	var sc model.ShoppingCart
	err := db.Where(`"UserID"=? and "GoodsID"=? and "SpecificationID"=?`, UserID, GoodsID, SpecificationID).First(&sc).Error
	log.Println(err)
	return sc
}
func (service ShoppingCartService) UpdateByUserIDAndID(db *gorm.DB, UserID types.PrimaryKey, GoodsID, SpecificationID types.PrimaryKey, Quantity uint) error {

	_sc := service.GetByGoodsIDAndSpecificationID(db, UserID, GoodsID, SpecificationID)
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
	return db.Model(&model.ShoppingCart{}).Where(`"GoodsID"=?`, GoodsID).Where(`"SpecificationID"=?`, SpecificationID).Where(`"UserID"=?`, UserID).Updates(_sc).Error
}

func (service ShoppingCartService) DeleteByUserIDAndGoodsIDAndSpecificationID(db *gorm.DB, UserID, GoodsID, SpecificationID types.PrimaryKey) error {
	return db.Where(&model.ShoppingCart{UserID: UserID, GoodsID: GoodsID, SpecificationID: SpecificationID}).Delete(&model.ShoppingCart{}).Error
}
func (service ShoppingCartService) DeleteListByIDs(db *gorm.DB, UserID types.PrimaryKey, GoodsID, SpecificationID types.PrimaryKey) error {
	return db.Where(&model.ShoppingCart{UserID: UserID, GoodsID: GoodsID, SpecificationID: SpecificationID}).Delete(&model.ShoppingCart{}).Error
}
func (service ShoppingCartService) FindShoppingCartByUserID(UserID types.PrimaryKey) []types.IEntity {
	Orm := singleton.Orm()
	//var list []model.ShoppingCart
	list := dao.Find(Orm, entity.ShoppingCart).Where(model.ShoppingCart{UserID: UserID}).List() //service.FindWhere(Orm, &list, model.ShoppingCart{UserID: UserID})
	return list
}
func (service ShoppingCartService) FindShoppingCartListDetails(UserID types.PrimaryKey) ([]AnalyseOrdersGoods, uint, error) {
	//Orm := Orm()
	ordersService := OrdersService{}

	list := service.FindShoppingCartByUserID(UserID)
	//results := make([]map[string]interface{}, 0)
	//var error error

	oredersGoodsList := make([]model.OrdersGoods, 0)

	for i := range list {
		oredersGoods := ordersService.createOrdersGoods(list[i].(*model.ShoppingCart))
		oredersGoodsList = append(oredersGoodsList, oredersGoods)
		//results[oredersGoods.OID]=append(results[oredersGoods.OID],oredersGoods)
	}
	return ordersService.AnalyseOrdersGoodsList(UserID, extends.Address{}, 0, oredersGoodsList)

}
