package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/service"
)

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
	pp := service.Payment.NewPaypal(ctx, oid)
	token, err := pp.GetAccessToken(ctx)
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
