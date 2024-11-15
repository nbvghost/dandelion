package serviceargument

import "time"

/**
* SUCCESS：支付成功
* REFUND：转入退款
* NOTPAY：未支付
* CLOSED：已关闭
* REVOKED：已撤销（仅付款码支付会返回）
* USERPAYING：用户支付中（仅付款码支付会返回）
* PAYERROR：支付失败（仅付款码支付会返回）
 */
type OrderQueryState string

const (
	OrderQueryState_SUCCESS    OrderQueryState = "SUCCESS"    //SUCCESS：支付成功
	OrderQueryState_REFUND     OrderQueryState = "REFUND"     //REFUND：转入退款
	OrderQueryState_NOTPAY     OrderQueryState = "NOTPAY"     //NOTPAY：未支付
	OrderQueryState_CLOSED     OrderQueryState = "CLOSED"     //CLOSED：已关闭
	OrderQueryState_REVOKED    OrderQueryState = "REVOKED"    //REVOKED：已撤销（仅付款码支付会返回）
	OrderQueryState_USERPAYING OrderQueryState = "USERPAYING" //USERPAYING：用户支付中（仅付款码支付会返回）
	OrderQueryState_PAYERROR   OrderQueryState = "PAYERROR"   //PAYERROR：支付失败（仅付款码支付会返回）
)

type OrderQueryResult struct {
	State            OrderQueryState
	PayerTotalAmount int64
	PayTime          time.Time
	OutTradeNo       string
	TransactionID    string
	Attach           string
}

type OrderResult struct {
	PrepayId string
}



