package internal

import (
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/service"
	"io"
	"log"
	"net/http"
	"time"
)

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
	pp:=service.Payment.NewPaypal(ctx,oid)
	token, err := pp.GetAccessToken()
	//token, err := generateAccessToken(ctx, oid)
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

	switch response.Status {
	case "CREATED":
	case "SAVED":
	case "APPROVED":
	case "VOIDED":
	case "COMPLETED":
	case "PAYER_ACTION_REQUIRED":
	}

	return responseData, nil
}