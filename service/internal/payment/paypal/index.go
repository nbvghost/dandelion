package paypal

import (
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/service/internal/configuration"
	"github.com/nbvghost/dandelion/service/serviceargument"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Service struct {
	OID     dao.PrimaryKey
	Context constrain.IContext
	configuration configuration.ConfigurationService
}

func (s *Service) CloseOrder(OrderNo string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Order(OrderNo string, title, description string, detail, openid string, IP string, Money uint, attach string) (*serviceargument.OrderResult, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) OrderQuery(OrderNo string) (*serviceargument.OrderQuery, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Refund(order *model.Orders, ordersGoods *model.OrdersGoods, reason string) error {
	panic("implement me")
}

func (s *Service) newAccessTokenRedisKey(oid dao.PrimaryKey) string {

	return fmt.Sprintf("payment:paypal:%d:access-token", oid)
}

// curl -v -X POST "https://api-m.sandbox.paypal.com/v1/oauth2/token"
// -u "<CLIENT_ID>:<CLIENT_SECRET>"
// -H "Content-Type: application/x-www-form-urlencoded"
// -d "grant_type=client_credentials"
type PaypalAccessToken struct {
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	AppId       string `json:"app_id"`
	ExpiresIn   int    `json:"expires_in"`
	Nonce       string `json:"nonce"`
}
type responseBody struct {
	PaypalAccessToken
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (s *Service) GetAccessToken() (*PaypalAccessToken, error) {

	configMap := s.configuration.GetConfigurations(s.OID, model.ConfigurationKeyPaymentPaypalClientId, model.ConfigurationKeyPaymentPaypalAppSecret)

	at := &responseBody{}
	clientId := configMap[model.ConfigurationKeyPaymentPaypalClientId]
	appSecret := configMap[model.ConfigurationKeyPaymentPaypalAppSecret]

	atJson, _ := s.Context.Redis().Get(s.Context, s.newAccessTokenRedisKey(s.OID))

	if len(atJson) > 0 {
		err := json.Unmarshal([]byte(atJson), &at.PaypalAccessToken)
		if err != nil {
			return nil, err
		}
		if len(at.PaypalAccessToken.AccessToken) > 0 {
			return &at.PaypalAccessToken, nil
		}

	}

	params := url.Values{}
	params.Set("grant_type", "client_credentials")

	client := &http.Client{}
	req, err := http.NewRequest("POST", key.BaseURL()+"/v1/oauth2/token", strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(clientId, appSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, at)
	if err != nil {
		return nil, err
	}

	err = s.Context.Redis().Set(s.Context, s.newAccessTokenRedisKey(s.OID), &at.PaypalAccessToken, time.Second*time.Duration(at.ExpiresIn))
	if err != nil {
		return nil, err
	}
	return &at.PaypalAccessToken, nil
}
