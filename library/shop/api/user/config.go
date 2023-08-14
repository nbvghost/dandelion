package user

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/configuration"
	"github.com/nbvghost/dandelion/service/order"
)

type Config struct {
	User                 *model.User `mapping:""`
	ConfigurationService configuration.ConfigurationService
	ShoppingCartService  order.ShoppingCartService
	Get                  struct {
	} `method:"Get"`
}

func (m *Config) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	defaultAddressList := dao.Find(db.Orm(), &model.Address{}).Where(`"UserID"=?`, context.UID()).Where(`"DefaultBilling"=true or "DefaultShipping"=true`).List()
	user := dao.GetByPrimaryKey(db.Orm(), &model.User{}, context.UID())
	userInfo := dao.GetBy(db.Orm(), &model.UserInfo{}, map[string]any{"UserID": context.UID()})
	data := m.ConfigurationService.GetConfigurations(m.User.OID, model.ConfigurationKeyAdvert, model.ConfigurationKeyPop, model.ConfigurationKeyQuickLink)
	ossUrl, err := oss.Url(context)
	if err != nil {
		return nil, err
	}
	return result.NewData(map[string]any{
		"Advert":         data[model.ConfigurationKeyAdvert],
		"Pop":            data[model.ConfigurationKeyPop],
		"QuickLink":      data[model.ConfigurationKeyQuickLink],
		"User":           user,
		"UserInfo":       userInfo,
		"DefaultAddress": defaultAddressList,
		"OSSUrl":         ossUrl,
	}), nil
}
