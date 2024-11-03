package goods

import (
	"fmt"
	"github.com/nbvghost/dandelion/domain/cache"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"gorm.io/gorm"
	"time"
)

func (m GoodsTypeService) CreateGoodsTypeUri(orm *gorm.DB, OID dao.PrimaryKey, text string) string {
	uri := cache.Cache.ChinesePinyinCache.AutoDetectUri(text)
	var gt model.GoodsType
	orm.Model(model.GoodsType{}).Where(map[string]interface{}{"OID": OID, "Uri": uri}).First(&gt)

	var gtc model.GoodsTypeChild
	orm.Model(model.GoodsTypeChild{}).Where(map[string]interface{}{"OID": OID, "Uri": uri}).First(&gtc)

	if gt.IsZero() && gtc.IsZero() {
		return uri
	}
	return fmt.Sprintf("%s-%s", uri, time.Now().Format("20060102150405"))
}

/*func (m GoodsTypeService) getGoodsTypeChildByUri(orm *gorm.DB, OID, GoodsTypeID dao.PrimaryKey, uri string) (model.GoodsTypeChild, error) {
	var gt model.GoodsTypeChild
	err := orm.Model(model.GoodsTypeChild{}).Where(map[string]interface{}{"OID": OID, "GoodsTypeID": GoodsTypeID, "Uri": uri}).First(&gt).Error
	return gt, err
}
*/
