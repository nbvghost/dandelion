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
	SessionAdmin          CookieKey = "ADMIN"         //商家后台
	SessionManager        CookieKey = "MANAGER"       //系统管理
	SessionUser           CookieKey = "USER"          //前台用户
	SessionOrganization   CookieKey = "Organization"  //前台用户
	SessionContentConfig  CookieKey = "ContentConfig" //
	SessionStore          CookieKey = "STORE"         //
	SessionAction         CookieKey = "ACTION"
	SessionUserID         CookieKey = "USERID"
	SessionOpenID         CookieKey = "OPENID"
	SessionMiniProgramKey CookieKey = "SESSIONMINIPROGRAMKEY"
	SessionConfirmOrders  CookieKey = "SESSIONCONFIRMORDERS"
	SessionRedirect       CookieKey = "REDIRECT"
	SessionCart           CookieKey = "CART"
	SessionCaptcha        CookieKey = "CAPTCHA"
	SessionSMSCode        CookieKey = "SMS_CODE"
)

var SessionGoodsViewFunc = func(goodsID uint) CookieKey {
	return CookieKey(fmt.Sprintf("SessionGoodsView_%v", goodsID))
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
	ScoreJournal_Type_Look_Article = 5
	ScoreJournal_Type_Share        = 6 //转发获历
	ScoreJournal_Type_InviteUser   = 7 //邀请好友
)
const (
	ContentTypeArticles = "articles"
)

var GWebSecretKey = encryption.SecretKey(encryption.Md5ByString("ds1f4ds524f52ds4f5ds4"))
