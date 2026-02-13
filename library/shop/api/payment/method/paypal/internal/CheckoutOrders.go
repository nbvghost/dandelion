package internal

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/service"
)

type CheckoutOrdersRequest struct {
	Intent        string               `json:"intent,omitempty"`
	PurchaseUnits []CheckoutOrdersUnit `json:"purchase_units,omitempty"`
	//Payer         CheckoutOrdersPayer         `json:"payer,omitempty"`
	//PaymentSource CheckoutOrdersPaymentSource `json:"payment_source,omitempty"`
}
type CheckoutOrdersResponse struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Id      string `json:"id"`
	Status  string `json:"status"`
	Links   []Link `json:"links"`
}

func CheckoutOrders(ctx constrain.IContext, oid dao.PrimaryKey, request *CheckoutOrdersRequest) (*CheckoutOrdersResponse, error) {
	pp := service.Payment.NewPaypal(ctx, oid)
	token, err := pp.GetAccessToken(ctx)
	//token, err := generateAccessToken(ctx, oid)
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
