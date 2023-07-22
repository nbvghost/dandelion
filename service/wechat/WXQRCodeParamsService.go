package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/library/db"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/tool/encryption"
	"github.com/nbvghost/tool/object"
)

type WXQRCodeParamsService struct {
	model.BaseDao
}

func (service WXQRCodeParamsService) addParams(key string, params map[string]interface{}) (*model.WXQRCodeParams, error) {
	b, _ := json.Marshal(params)
	wxParams := &model.WXQRCodeParams{}
	wxParams.CodeKey = key
	wxParams.Params = string(b)
	err := dao.Create(db.Orm(), wxParams)
	if err != nil {
		return wxParams, err
	}
	return wxParams, nil
}
func (service WXQRCodeParamsService) getParams(CodeKey string) (*model.WXQRCodeParams, error) {
	wxParams := &model.WXQRCodeParams{}
	db.Orm().Model(&model.WXQRCodeParams{}).Where(`"CodeKey"=?`, CodeKey).First(wxParams)
	if wxParams.ID == 0 {
		return wxParams, errors.New("NOT FOUND")
	}
	return wxParams, nil
}
func (service WXQRCodeParamsService) EncodeShareKey(UserID dao.PrimaryKey, ProductID uint) string {

	key := encryption.Md5ByString(fmt.Sprintf("%v%v", UserID, ProductID))
	wxParams, err := service.getParams(key)
	if err != nil {
		wxParams, _ := service.addParams(key, map[string]interface{}{
			"UserID":    UserID,
			"ProductID": ProductID,
		})
		return wxParams.CodeKey
	} else {
		return wxParams.CodeKey
	}
}
func (service WXQRCodeParamsService) DecodeShareKey(ShareKey string) (UserID dao.PrimaryKey, ProductID uint) {
	wxParams, _ := service.getParams(ShareKey)

	paramsMap := make(map[string]interface{})

	json.Unmarshal([]byte(wxParams.Params), &paramsMap)

	UserID = dao.PrimaryKey(object.ParseInt(paramsMap["UserID"]))
	ProductID = uint(object.ParseInt(paramsMap["ProductID"]))
	return
}
