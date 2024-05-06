package serviceargument

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"time"
)

type ListQueryParam struct {
	GoodsTypeID      dao.PrimaryKey
	GoodsTypeChildID dao.PrimaryKey
}

type ListOrdersQueryParam struct {
	UserID             dao.PrimaryKey
	Status             []model.OrdersStatus
	StartDate, EndDate time.Time
}

type RefundNotifyData struct {
	Mchid         string    `json:"mchid"`
	OutTradeNo    string    `json:"out_trade_no"`
	TransactionId string    `json:"transaction_id"`
	OutRefundNo   string    `json:"out_refund_no"`
	RefundId      string    `json:"refund_id"`
	RefundStatus  string    `json:"refund_status"`
	SuccessTime   time.Time `json:"success_time"`
	Amount        struct {
		Total       int `json:"total"`
		Refund      int `json:"refund"`
		PayerTotal  int `json:"payer_total"`
		PayerRefund int `json:"payer_refund"`
	} `json:"amount"`
	UserReceivedAccount string `json:"user_received_account"`
}