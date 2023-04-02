package user

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/configuration"
)

type Config struct {
	User                 *model.User `mapping:""`
	ConfigurationService configuration.ConfigurationService
	Get                  struct {
		AppId string `uri:"AppId"`
	} `method:"Get"`
}

func (m *Config) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	defaultAddressList := dao.Find(singleton.Orm(), &model.Address{}).Where(`"UserID"=?`, context.UID()).Where(`"DefaultBilling"=true or "DefaultShipping"=true`).List()

	user := dao.GetByPrimaryKey(singleton.Orm(), &model.User{}, context.UID())

	userInfo := dao.GetBy(singleton.Orm(), &model.UserInfo{}, map[string]any{"UserID": context.UID()})

	data := m.ConfigurationService.GetConfigurations(m.User.OID, model.ConfigurationKeyAdvert, model.ConfigurationKeyPop, model.ConfigurationKeyQuickLink)
	return result.NewData(map[string]any{
		"Advert":         data[model.ConfigurationKeyAdvert],
		"Pop":            data[model.ConfigurationKeyPop],
		"QuickLink":      data[model.ConfigurationKeyQuickLink],
		"User":           user,
		"UserInfo":       userInfo,
		"DefaultAddress": defaultAddressList,
	}), nil
}
