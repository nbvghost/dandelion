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

func (m *Config) Handle(ctx constrain.IContext) (r constrain.IResult, err error) {
	defaultAddressList := dao.Find(db.GetDB(ctx), &model.Address{}).Where(`"UserID"=?`, ctx.UID()).Where(`"DefaultBilling"=true or "DefaultShipping"=true`).List()
	user := dao.GetByPrimaryKey(db.GetDB(ctx), &model.User{}, ctx.UID())

	userInfo := service.User.GetUserInfo(ctx, ctx.UID())

	data := service.Configuration.GetConfigurations(ctx,
		m.User.OID,
		model.ConfigurationKeyAdvert,
		model.ConfigurationKeyPop,
		model.ConfigurationKeyQuickLink,
		model.ConfigurationKeyPaymentPaypalClientId,
	)
	ossUrl, err := oss.Url(ctx)
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
