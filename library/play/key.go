package play

import (
	"fmt"

	"github.com/nbvghost/gweb"

	"github.com/nbvghost/tool/encryption"
)

type CookieKey string

const CookieKeyManager CookieKey = "dandelion_manager"
const CookieKeyAdmin CookieKey = "dandelion_admin"
const CookieKeyUser CookieKey = "dandelion_user"

const (
	SessionAdmin          gweb.AttributesKey = "ADMIN"         //商家后台
	SessionManager        gweb.AttributesKey = "MANAGER"       //系统管理
	SessionUser           gweb.AttributesKey = "USER"          //前台用户
	SessionOrganization   gweb.AttributesKey = "Organization"  //前台用户
	SessionContentConfig  gweb.AttributesKey = "ContentConfig" //
	SessionStore          gweb.AttributesKey = "STORE"         //
	SessionAction         gweb.AttributesKey = "ACTION"
	SessionUserID         gweb.AttributesKey = "USERID"
	SessionOpenID         gweb.AttributesKey = "OPENID"
	SessionMiniProgramKey gweb.AttributesKey = "SESSIONMINIPROGRAMKEY"
	SessionConfirmOrders  gweb.AttributesKey = "SESSIONCONFIRMORDERS"
	SessionRedirect       gweb.AttributesKey = "REDIRECT"
	SessionCart           gweb.AttributesKey = "CART"
	SessionCaptcha        gweb.AttributesKey = "CAPTCHA"
	SessionSMSCode        gweb.AttributesKey = "SMS_CODE"
)

var SessionGoodsViewFunc = func(goodsID uint) gweb.AttributesKey {
	return gweb.AttributesKey(fmt.Sprintf("SessionGoodsView_%v", goodsID))
}

const (
	QRCodeCreateType_Article = "A" //文章二维码
)
const (
	ActionKey_add    string = "add"
	ActionKey_save   string = "save"
	ActionKey_change string = "change"
	ActionKey_get    string = "get"
	ActionKey_one    string = "one"
	ActionKey_list   string = "list"
	ActionKey_del    string = "del"
)
const (
	WxConfigType_miniprogram = "miniprogram" //小程序
	WxConfigType_miniweb     = "miniweb"     //公众号
)

const (
	OS_Order    = "Order"    // order=下单成功，待付款
	OS_Pay      = "Pay"      // pay=支付成功，待发货
	OS_Deliver  = "Deliver"  // deliver=发货成功，待收货
	OS_Refund   = "Refund"   // Refund=订单退款退货中->所有子商品状态为空或OGRefundOK->返回Deliver状态
	OS_RefundOk = "RefundOk" // Orders 下的所有ordergoods 全部退款，orders 改为 RefundOk
	OS_OrderOk  = "OrderOk"  // order_ok=订单确认完成
	OS_Cancel   = "Cancel"   //订单等待取消
	OS_CancelOk = "CancelOk" //订单已经取消
	OS_Delete   = "Delete"   // delete=删除
	OS_Closed   = "Closed"   // 已经关闭
)
const (
	//OrdersGoods,Voucher,ScoreGoods
	CardItem_Type_OrdersGoods = "OrdersGoods"
	CardItem_Type_Voucher     = "Voucher"
	CardItem_Type_ScoreGoods  = "ScoreGoods"
)
const (
	OrdersType_Goods        = "Goods"        //商品购买订单
	OrdersType_GoodsPackage = "GoodsPackage" //合并下单
	OrdersType_Supply       = "Supply"       //充值
)
const (
	SupplyType_Store = "Store" //门店充值
	SupplyType_User  = "User"  //普通用户充值
)
const (
	Paging int = 10
)
const (
	OS_OGAskRefund      = "OGAskRefund"      // OGAskRefund=申请，申请退货退款
	OS_OGRefundNo       = "OGRefundNo"       // OGRefundOK=拒绝子商品，确认退货款
	OS_OGRefundOk       = "OGRefundOk"       // OGRefundOK=允许子商品，确认退货款
	OS_OGRefundInfo     = "OGRefundInfo"     // OGRefundInfo=用户填写信息，允许子商品，确认退货款
	OS_OGRefundComplete = "OGRefundComplete" // OGRefund=完成子商品，用户邮寄商品，商家待收货
)

const (
	StoreJournal_Type_ZZHX = 1 //自主核销
	StoreJournal_Type_CZ   = 2 //在线充值
	StoreJournal_Type_HX   = 3 //来自用户订单的商品核销
	StoreJournal_Type_SG   = 4 //ScoreGoods 积分商品核销
	StoreJournal_Type_FL   = 5 //Voucher 福利卷核销
	StoreJournal_Type_TX   = 6 //店员提现
)
const (
	UserJournal_Type_HX        = 1 //核销
	UserJournal_Type_LEVE      = 2 //下单，上下级结算佣金
	UserJournal_Type_TX        = 3 //提现
	UserJournal_Type_USER_LEVE = 4 //成为上下级，结算佣金
)
const (
	OrganizationJournal_Goods     = 1 //商品销售
	OrganizationJournal_Brokerage = 2 //商品销售用户的佣金
)
const (
	ScoreJournal_Type_GM           = 1 //购买商品
	ScoreJournal_Type_DH           = 2 //积分兑换商品
	ScoreJournal_Type_LEVE         = 3 //上下级结算佣金,获取的积分
	ScoreJournal_Type_DaySign      = 4 //签到送积分
	ScoreJournal_Type_Look_Article = 5 //看文章
	ScoreJournal_Type_Share        = 6 //转发获历
	ScoreJournal_Type_InviteUser   = 7 //邀请好友
)
const (
	ContentTypeArticles = "articles"
)

var GWebSecretKey = encryption.SecretKey(encryption.Md5ByString("ds1f4ds524f52ds4f5ds4"))
