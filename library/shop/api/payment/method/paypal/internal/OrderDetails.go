package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/service"
)

type OrderDetailsResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Links  []Link `json:"links"`
}

func OrderDetails(ctx constrain.IContext, oid dao.PrimaryKey, id string) (*OrderDetailsResponse, error) {
	pp := service.Payment.NewPaypal(ctx, oid)
	token, err := pp.GetAccessToken(ctx)

	//token, err := generateAccessToken(ctx, oid)
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
