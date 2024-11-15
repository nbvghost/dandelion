package play

import (
	"fmt"
	"github.com/nbvghost/tool/encryption"
)

type CookieKey string

const CookieKeyManager CookieKey = "dandelion_manager"
const CookieKeyAdmin CookieKey = "dandelion_admin"
const CookieKeyUser CookieKey = "dandelion_user"

const (
	SessionAdmin          AttributesKey = "ADMIN"         //商家后台
	SessionManager        AttributesKey = "MANAGER"       //系统管理
	SessionUser           AttributesKey = "USER"          //前台用户
	SessionOrganization   AttributesKey = "Organization"  //前台用户
	SessionContentConfig  AttributesKey = "ContentConfig" //
	SessionStore          AttributesKey = "STORE"         //
	SessionAction         AttributesKey = "ACTION"
	SessionUserID         AttributesKey = "USERID"
	SessionOpenID         AttributesKey = "OPENID"
	SessionMiniProgramKey AttributesKey = "SESSIONMINIPROGRAMKEY"
	SessionConfirmOrders  AttributesKey = "SESSIONCONFIRMORDERS"
	SessionRedirect       AttributesKey = "REDIRECT"
	SessionCart           AttributesKey = "CART"
	SessionCaptcha        AttributesKey = "CAPTCHA"
	SessionSMSCode        AttributesKey = "SMS_CODE"
)

var SessionGoodsViewFunc = func(goodsID uint) AttributesKey {
	return AttributesKey(fmt.Sprintf("SessionGoodsView_%v", goodsID))
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
	//OrdersGoods,Voucher,ScoreGoods
	CardItem_Type_OrdersGoods = "OrdersGoods"
	CardItem_Type_Voucher     = "Voucher"
	CardItem_Type_ScoreGoods  = "ScoreGoods"
)


const (
	SupplyType_Store = "Store" //门店充值
	SupplyType_User  = "User"  //普通用户充值
)
const (
	Paging int = 10
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
	OrganizationJournal_Goods     = 1 //商品销售
	OrganizationJournal_Brokerage = 2 //商品销售用户的佣金
)

const (
	ContentTypeArticles = "articles"
)

var GWebSecretKey = encryption.SecretKey(encryption.Md5ByString("ds1f4ds524f52ds4f5ds4"))
