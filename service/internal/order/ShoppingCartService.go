package order

import (
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service/internal/activity"
	"github.com/nbvghost/dandelion/service/internal/company"
	"github.com/nbvghost/dandelion/service/internal/goods"
	"log"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
)

type ShoppingCartService struct {
	model.BaseDao
	Goods        goods.GoodsService
	TimeSell     activity.TimeSellService
	Organization company.OrganizationService
}

func (service ShoppingCartService) GetGSIDs(UserID dao.PrimaryKey, GSIDs []string) []model.ShoppingCart {
	Orm := db.Orm()
	var scs []model.ShoppingCart
	//Orm := Orm()
	//return Orm.Where("GSID in (?)", IDs).Where(&model.ShoppingCart{UserID: UserID}).Delete(&model.ShoppingCart{}).Error
	err := Orm.Where(`"GSID" in (?) and "UserID"=?`, GSIDs, UserID).Find(&scs).Error
	log.Println(err)
	return scs
}
func (service ShoppingCartService) GetByGoodsIDAndSpecificationID(db *gorm.DB, UserID dao.PrimaryKey, GoodsID, SpecificationID dao.PrimaryKey) model.ShoppingCart {
	var sc model.ShoppingCart
	err := db.Where(`"UserID"=? and "GoodsID"=? and "SpecificationID"=?`, UserID, GoodsID, SpecificationID).First(&sc).Error
	log.Println(err)
	return sc
}
func (service ShoppingCartService) UpdateByUserIDAndID(db *gorm.DB, UserID dao.PrimaryKey, GoodsID, SpecificationID dao.PrimaryKey, Quantity uint) error {

	_sc := service.GetByGoodsIDAndSpecificationID(db, UserID, GoodsID, SpecificationID)
	//err := service.Get(Orm, ID, &_sc)
	/*err := service.GetGSID(UserID,GSID)
	if err != nil {
		return err
	}*/

	//var specification model.Specification
	//util.JSONToStruct(_sc.Specification, &specification)
	/*err = service.Goods.GetSpecification(_sc.SpecificationID, &specification)
	if err != nil {
		return err
	}*/
	specification := dao.GetByPrimaryKey(db, &model.Specification{}, SpecificationID).(*model.Specification)
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

func (service ShoppingCartService) DeleteByUserIDAndGoodsIDAndSpecificationID(db *gorm.DB, UserID, GoodsID, SpecificationID dao.PrimaryKey) error {
	return db.Where(&model.ShoppingCart{UserID: UserID, GoodsID: GoodsID, SpecificationID: SpecificationID}).Delete(&model.ShoppingCart{}).Error
}
func (service ShoppingCartService) DeleteListByIDs(db *gorm.DB, UserID dao.PrimaryKey, GoodsID, SpecificationID dao.PrimaryKey) error {
	return db.Where(&model.ShoppingCart{UserID: UserID, GoodsID: GoodsID, SpecificationID: SpecificationID}).Delete(&model.ShoppingCart{}).Error
}
func (service ShoppingCartService) FindShoppingCartByUserID(UserID dao.PrimaryKey) []dao.IEntity {
	Orm := db.Orm()
	//var list []model.ShoppingCart
	list := dao.Find(Orm, entity.ShoppingCart).Where(model.ShoppingCart{UserID: UserID}).List() //service.FindWhere(Orm, &list, model.ShoppingCart{UserID: UserID})
	return list
}
func (service ShoppingCartService) FindShoppingCartListCount(UserID dao.PrimaryKey) (uint, error) {
	list := service.FindShoppingCartByUserID(UserID)
	return uint(len(list)), nil
}
