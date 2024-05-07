package paypal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/service/internal/configuration"
	"github.com/nbvghost/dandelion/service/serviceargument"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Service struct {
	OID           dao.PrimaryKey
	Context       constrain.IContext
	configuration configuration.ConfigurationService
}

func (m *Service) CloseOrder(OrderNo string) error {
	//TODO implement me
	panic("implement me")
}

func (m *Service) Order(OrderNo string, title, description string, detail, openid string, IP string, Money uint, attach string) (*serviceargument.OrderResult, error) {
	//TODO implement me
	panic("implement me")
}

func (m *Service) OrderQuery(OrderNo string) (*serviceargument.OrderQuery, error) {
	//TODO implement me
	panic("implement me")
}

type RefundRequest struct {
	/*Amount struct {
		Value        string `json:"value"`
		CurrencyCode string `json:"currency_code"`
	} `json:"amount"`*/
	CustomId string `json:"custom_id"`
	//InvoiceId          string `json:"invoice_id"`
	NoteToPayer string `json:"note_to_payer"`
	/*PaymentInstruction struct {
		PlatformFees []struct {
			Amount struct {
				CurrencyCode string `json:"currency_code"`
				Value        string `json:"value"`
			} `json:"amount"`
		} `json:"platform_fees"`
	} `json:"payment_instruction"`*/
}
type RefundResponse struct {
	Status        string `json:"status"`
	StatusDetails struct {
		Reason string `json:"reason"`
	} `json:"status_details"`
}

func (m *Service) Refund(order *model.Orders, ordersGoods *model.OrdersGoods, reason string) error {

	token, err := m.GetAccessToken()
	//token, err := generateAccessToken(ctx, oid)
	if err != nil {
		return err
	}

	body, err := json.Marshal(&RefundRequest{
		CustomId:    order.OrderNo,
		NoteToPayer: reason,
	})
	if err != nil {
		return err
	}

	client := &http.Client{}
	///v2/checkout/orders/${orderId}/capture
	//req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/payments/captures/%s/refund", key.BaseURL(), order.TransactionID), bytes.NewReader(body))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/payments/captures/%s/refund", key.BaseURL(), order.TransactionID), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("PayPal-Request-Id", uuid.New().String())

	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err = io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	log.Println(string(body))
	responseData := &RefundResponse{}
	err = json.Unmarshal(body, responseData)
	if err != nil {
		return err
	}
	switch responseData.Status {
	case "CANCELLED":
		return errors.New("The refund was cancelled")
	case "FAILED":
		return errors.New(responseData.StatusDetails.Reason)
	case "PENDING":
		return errors.New(responseData.StatusDetails.Reason)
	case "COMPLETED":
	default:
		return errors.New("refund error")
	}
	return nil
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

func (m *Service) GetAccessToken() (*PaypalAccessToken, error) {

	configMap := m.configuration.GetConfigurations(m.OID, model.ConfigurationKeyPaymentPaypalClientId, model.ConfigurationKeyPaymentPaypalAppSecret)

	at := &responseBody{}
	clientId := configMap[model.ConfigurationKeyPaymentPaypalClientId]
	appSecret := configMap[model.ConfigurationKeyPaymentPaypalAppSecret]

	atJson, _ := m.Context.Redis().Get(m.Context, key.NewPaypalAccessTokenRedisKey(m.OID))

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

	err = m.Context.Redis().Set(m.Context, key.NewPaypalAccessTokenRedisKey(m.OID), &at.PaypalAccessToken, time.Second*time.Duration(at.ExpiresIn))
	if err != nil {
		return nil, err
	}
	return &at.PaypalAccessToken, nil
}
