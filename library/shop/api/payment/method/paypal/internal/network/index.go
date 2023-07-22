package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/service/configuration"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Amount struct {
	CurrencyCode string `json:"currency_code,omitempty"`
	Value        string `json:"value,omitempty"`
}
type Name struct {
	GivenName string `json:"given_name,omitempty"` //名
	Surname   string `json:"surname,omitempty"`    //姓
	FullName  string `json:"full_name,omitempty"`
}

func (m *Name) GetFullName() string {
	return m.GivenName + " " + m.Surname
}

type Shipping struct {
	Name    *Name    `json:"name,omitempty"`
	Type    string   `json:"type,omitempty"` //SHIPPING
	Address *Address `json:"address,omitempty"`
}
type CheckoutOrdersUnit struct {
	ReferenceId string    `json:"reference_id,omitempty"`
	Description string    `json:"description,omitempty"`
	Amount      Amount    `json:"amount,omitempty"`
	Shipping    *Shipping `json:"shipping,omitempty"`
}
type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type Payer struct {
	PayerId      string `json:"payer_id,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	Name         Name   `json:"name,omitempty"`
	Phone        struct {
		PhoneType   string `json:"phone_type,omitempty"`
		PhoneNumber struct {
			NationalNumber string `json:"national_number,omitempty"`
		} `json:"phone_number,omitempty"`
	} `json:"phone,omitempty"`
	BirthDate string `json:"birth_date,omitempty"`
	TaxInfo   struct {
		TaxId     string `json:"tax_id,omitempty"`
		TaxIdType string `json:"tax_id_type,omitempty"`
	} `json:"tax_info,omitempty"`
	Address Address `json:"address,omitempty"`
}
type Address struct {
	AddressLine1 string `json:"address_line_1,omitempty"`
	AddressLine2 string `json:"address_line_2,omitempty"`
	AdminArea2   string `json:"admin_area_2,omitempty"`
	AdminArea1   string `json:"admin_area_1,omitempty"`
	PostalCode   string `json:"postal_code,omitempty"`
	CountryCode  string `json:"country_code,omitempty"`
}

func (m *Address) SetAddress(address *model.Address) *Address {
	if address == nil {
		return nil
	}
	m.AddressLine1 = address.Detail
	if len(address.Company) > 0 {
		m.AddressLine2 = fmt.Sprintf("(%s)", address.Company)
	}
	m.AdminArea1 = address.CountyName + "." + address.ProvinceName
	m.AdminArea2 = address.CityName
	m.PostalCode = address.PostalCode
	m.CountryCode = address.CountyCode
	return m
}

type CheckoutOrdersCard struct {
	Name           string  `json:"name,omitempty"`
	Number         string  `json:"number,omitempty"`
	SecurityCode   string  `json:"security_code,omitempty"`
	Expiry         string  `json:"expiry,omitempty"`
	BillingAddress Address `json:"billing_address,omitempty"`
}
type CheckoutOrdersPaymentSource struct {
	Card CheckoutOrdersCard `json:"card,omitempty"`
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

func newAccessTokenRedisKey(oid dao.PrimaryKey) string {

	return fmt.Sprintf("payment:paypal:%d:access-token", oid)
}
func generateAccessToken(ctx constrain.IContext, oid dao.PrimaryKey) (*PaypalAccessToken, error) {
	configurationService := configuration.ConfigurationService{}
	configMap := configurationService.GetConfigurations(oid, model.ConfigurationKeyPaymentPaypalClientId, model.ConfigurationKeyPaymentPaypalAppSecret)

	at := &responseBody{}
	clientId := configMap[model.ConfigurationKeyPaymentPaypalClientId]
	appSecret := configMap[model.ConfigurationKeyPaymentPaypalAppSecret]

	atJson, _ := ctx.Redis().Get(ctx, newAccessTokenRedisKey(oid))

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

	err = ctx.Redis().Set(ctx, newAccessTokenRedisKey(oid), &at.PaypalAccessToken, time.Second*time.Duration(at.ExpiresIn))
	if err != nil {
		return nil, err
	}
	return &at.PaypalAccessToken, nil
}

type OrderDetailsResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Links  []Link `json:"links"`
}

func OrderDetails(ctx constrain.IContext, oid dao.PrimaryKey, id string) (*OrderDetailsResponse, error) {
	token, err := generateAccessToken(ctx, oid)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/checkout/orders/%s", key.BaseURL(), id), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("PayPal-Request-Id", uuid.New().String())

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	log.Println(string(body))
	responseData := &OrderDetailsResponse{}
	err = json.Unmarshal(body, responseData)
	if err != nil {
		return nil, err
	}
	return responseData, nil
}

type UpdateOrderChange struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value any    `json:"value"`
}
type UpdateOrderRequest struct {
	Id         string
	ChangeList []UpdateOrderChange
}

func UpdateOrder(ctx constrain.IContext, oid dao.PrimaryKey, request *UpdateOrderRequest) error {
	token, err := generateAccessToken(ctx, oid)
	if err != nil {
		return err
	}

	body, err := json.Marshal(request.ChangeList)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/v2/checkout/orders/%s", key.BaseURL(), request.Id), strings.NewReader(string(body)))
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

	if response.StatusCode != 204 {
		return errors.New("order modification failed")
	}
	return nil
}

type CheckoutOrdersRequest struct {
	Intent        string               `json:"intent,omitempty"`
	PurchaseUnits []CheckoutOrdersUnit `json:"purchase_units,omitempty"`
	//Payer         CheckoutOrdersPayer         `json:"payer,omitempty"`
	//PaymentSource CheckoutOrdersPaymentSource `json:"payment_source,omitempty"`
}
type CheckoutOrdersResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Links  []Link `json:"links"`
}

func CheckoutOrders(ctx constrain.IContext, oid dao.PrimaryKey, request *CheckoutOrdersRequest) (*CheckoutOrdersResponse, error) {
	token, err := generateAccessToken(ctx, oid)
	if err != nil {
		return nil, err
	}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", key.BaseURL()+"/v2/checkout/orders", strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("PayPal-Request-Id", uuid.New().String())

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	log.Println(string(body))
	responseData := &CheckoutOrdersResponse{}
	err = json.Unmarshal(body, responseData)
	if err != nil {
		return nil, err
	}
	return responseData, nil
}

type CaptureRequest struct {
	PaypalOrderID string
}
type CaptureResponse struct {
	Id            string `json:"id"`
	Status        string `json:"status"`
	PaymentSource struct {
		Paypal struct {
			EmailAddress string  `json:"email_address"`
			AccountId    string  `json:"account_id"`
			Name         Name    `json:"name"`
			Address      Address `json:"address"`
		} `json:"paypal"`
	} `json:"payment_source"`
	PurchaseUnits []struct {
		ReferenceId string   `json:"reference_id"`
		Shipping    Shipping `json:"shipping"`
		Payments    struct {
			Captures []struct {
				Id               string `json:"id"`
				Status           string `json:"status"`
				Amount           Amount `json:"amount"`
				FinalCapture     bool   `json:"final_capture"`
				SellerProtection struct {
					Status            string   `json:"status"`
					DisputeCategories []string `json:"dispute_categories"`
				} `json:"-"`
				SellerReceivableBreakdown struct {
					GrossAmount Amount `json:"gross_amount"`
					NetAmount   Amount `json:"net_amount"`
					PaypalFee   Amount `json:"paypal_fee"`
				} `json:"-"`
				Links      []Link    `json:"links"`
				CreateTime time.Time `json:"create_time"`
				UpdateTime time.Time `json:"update_time"`
			} `json:"captures"`
		} `json:"payments"`
	} `json:"purchase_units"`
	Payer Payer  `json:"payer"`
	Links []Link `json:"links"`
}

func Capture(ctx constrain.IContext, oid dao.PrimaryKey, request *CaptureRequest) (*CaptureResponse, error) {
	token, err := generateAccessToken(ctx, oid)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	///v2/checkout/orders/${orderId}/capture
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/checkout/orders/%s/capture", key.BaseURL(), request.PaypalOrderID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("PayPal-Request-Id", uuid.New().String())

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	log.Println(string(body))
	responseData := &CaptureResponse{}
	err = json.Unmarshal(body, responseData)
	if err != nil {
		return nil, err
	}
	return responseData, nil
}
