package account

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type Config struct {
	User *model.User `mapping:""`
	Get  struct{}    `method:"Get"`
}

func (m *Config) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	defaultAddressList := dao.Find(db.Orm(), &model.Address{}).Where(`"UserID"=?`, context.UID()).Where(`"DefaultBilling"=true or "DefaultShipping"=true`).List()
	user := dao.GetByPrimaryKey(db.Orm(), &model.User{}, context.UID())

	userInfo := service.User.GetUserInfo(context.UID())

	data := service.Configuration.GetConfigurations(
		m.User.OID,
		model.ConfigurationKeyAdvert,
		model.ConfigurationKeyPop,
		model.ConfigurationKeyQuickLink,
		model.ConfigurationKeyPaymentPaypalClientId,
	)
	ossUrl, err := oss.Url(context)
	if err != nil {
		return nil, err
	}
	return result.NewData(map[string]any{
		"Advert":                data[model.ConfigurationKeyAdvert],
		"Pop":                   data[model.ConfigurationKeyPop],
		"QuickLink":             data[model.ConfigurationKeyQuickLink],
		"PaymentPaypalClientId": data[model.ConfigurationKeyPaymentPaypalClientId],
		"User":                  user,
		"UserInfo": map[string]any{
			"AllowAssistance": userInfo.GetAllowAssistance(),
			"Subscribe":       userInfo.GetSubscribe(),
		},
		"DefaultAddress": defaultAddressList,
		"OSSUrl":         ossUrl,
	}), nil
}
