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
	"github.com/nbvghost/tool/object"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Service struct {
	OID           dao.PrimaryKey
	Context       constrain.IServiceContext
	configuration configuration.ConfigurationService
}

func (m *Service) Deliver(orders *model.Orders) error {

	return nil
}

func (m *Service) CloseOrder(OrderNo string) error {
	//TODO implement me
	panic("implement me")
}

func (m *Service) Order(OrderNo string, title, description string, detail, openid string, IP string, Money uint, ordersType model.OrdersType) (*serviceargument.OrderResult, error) {
	//TODO implement me
	panic("implement me")
}

type OrderQueryResponse struct {
	Status        Status                 `json:"status"`
	StatusDetails StatusDetails          `json:"status_details"`
	Amount        serviceargument.Amount `json:"amount"`
	CreateTime    string                 `json:"create_time"`
}

func (m *Service) OrderQuery(orders *model.Orders) (*serviceargument.OrderQueryResult, error) {
	token, err := m.GetAccessToken()
	//token, err := generateAccessToken(ctx, oid)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	///v2/checkout/orders/${orderId}/capture
	//req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/payments/captures/%s/refund", key.BaseURL(), order.TransactionID), bytes.NewReader(body))
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/payments/captures/%s", key.BaseURL(), orders.TransactionID), nil)
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
	responseData := &OrderQueryResponse{}
	err = json.Unmarshal(body, responseData)
	if err != nil {
		return nil, err
	}

	createTime, err := time.ParseInLocation(time.RFC3339, responseData.CreateTime, time.Local)
	if err != nil {
		return nil, err
	}
	oqr := &serviceargument.OrderQueryResult{
		PayerTotalAmount: int64(object.ParseInt(object.ParseFloat(responseData.Amount.Value) * 100)),
		PayTime:          createTime,
		OutTradeNo:       orders.OrderNo,
		TransactionID:    orders.TransactionID,
		Attach:           "",
	}

	/*	const (
		OrderQueryState_SUCCESS    OrderQueryState = "SUCCESS"    //SUCCESS：支付成功
		OrderQueryState_REFUND     OrderQueryState = "REFUND"     //REFUND：转入退款
		OrderQueryState_NOTPAY     OrderQueryState = "NOTPAY"     //NOTPAY：未支付
		OrderQueryState_CLOSED     OrderQueryState = "CLOSED"     //CLOSED：已关闭
		OrderQueryState_REVOKED    OrderQueryState = "REVOKED"    //REVOKED：已撤销（仅付款码支付会返回）
		OrderQueryState_USERPAYING OrderQueryState = "USERPAYING" //USERPAYING：用户支付中（仅付款码支付会返回）
		OrderQueryState_PAYERROR   OrderQueryState = "PAYERROR"   //PAYERROR：支付失败（仅付款码支付会返回）
	)*/

	switch responseData.Status {
	case Status_COMPLETED:
		oqr.State = serviceargument.OrderQueryState_SUCCESS
	case Status_DECLINED:
		oqr.State = serviceargument.OrderQueryState_NOTPAY
	case Status_PARTIALLY_REFUNDED:
		oqr.State = serviceargument.OrderQueryState_REFUND
	case Status_PENDING:
		oqr.State = serviceargument.OrderQueryState_USERPAYING
	case Status_REFUNDED:
		oqr.State = serviceargument.OrderQueryState_REFUND
	case Status_FAILED:
		oqr.State = serviceargument.OrderQueryState_PAYERROR
	default:
		return nil, errors.New(fmt.Sprintf("error status:%s", responseData.Status))
	}
	return oqr, nil
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
type Status string

const (
	Status_COMPLETED          Status = "COMPLETED"          // The funds for this captured payment were credited to the payee's PayPal account.
	Status_DECLINED           Status = "DECLINED"           // The funds could not be captured.
	Status_PARTIALLY_REFUNDED Status = "PARTIALLY_REFUNDED" // An amount less than this captured payment's amount was partially refunded to the payer.
	Status_PENDING            Status = "PENDING"            // The funds for this captured payment was not yet credited to the payee's PayPal account. For more information, see status.details.
	Status_REFUNDED           Status = "REFUNDED"           // An amount greater than or equal to this captured payment's amount was refunded to the payer.
	Status_FAILED             Status = "FAILED"             // There was an error while capturing payment.
)

type StatusDetailsReason string

const (
	StatusDetailsReason_BUYER_COMPLAINT                             StatusDetailsReason = "BUYER_COMPLAINT"                             // The payer initiated a dispute for this captured payment with PayPal.
	StatusDetailsReason_CHARGEBACK                                  StatusDetailsReason = "CHARGEBACK"                                  // The captured funds were reversed in response to the payer disputing this captured payment with the issuer of the financial instrument used to pay for this captured payment.
	StatusDetailsReason_ECHECK                                      StatusDetailsReason = "ECHECK"                                      // The payer paid by an eCheck that has not yet cleared.
	StatusDetailsReason_INTERNATIONAL_WITHDRAWAL                    StatusDetailsReason = "INTERNATIONAL_WITHDRAWAL"                    // Visit your online account. In your Account Overview, accept and deny this payment.
	StatusDetailsReason_OTHER                                       StatusDetailsReason = "OTHER"                                       // No additional specific reason can be provided. For more information about this captured payment, visit your account online or contact PayPal.
	StatusDetailsReason_PENDING_REVIEW                              StatusDetailsReason = "PENDING_REVIEW"                              // The captured payment is pending manual review.
	StatusDetailsReason_RECEIVING_PREFERENCE_MANDATES_MANUAL_ACTION StatusDetailsReason = "RECEIVING_PREFERENCE_MANDATES_MANUAL_ACTION" // The payee has not yet set up appropriate receiving preferences for their account. For more information about how to accept or deny this payment, visit your account online. This reason is typically offered in scenarios such as when the currency of the captured payment is different from the primary holding currency of the payee.
	StatusDetailsReason_REFUNDED                                    StatusDetailsReason = "REFUNDED"                                    // The captured funds were refunded.
	StatusDetailsReason_TRANSACTION_APPROVED_AWAITING_FUNDING       StatusDetailsReason = "TRANSACTION_APPROVED_AWAITING_FUNDING"       // The payer must send the funds for this captured payment. This code generally appears for manual EFTs.
	StatusDetailsReason_UNILATERAL                                  StatusDetailsReason = "UNILATERAL"                                  // The payee does not have a PayPal account.
	StatusDetailsReason_VERIFICATION_REQUIRED                       StatusDetailsReason = "VERIFICATION_REQUIRED"                       // The payee's PayPal account is not verified.
)

type StatusDetails struct {
	Reason StatusDetailsReason `json:"reason"`
}
type RefundResponse struct {
	Status        Status        `json:"status"`
	StatusDetails StatusDetails `json:"status_details"`
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
		return errors.New(string(responseData.StatusDetails.Reason))
	case "PENDING":
		return errors.New(string(responseData.StatusDetails.Reason))
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
