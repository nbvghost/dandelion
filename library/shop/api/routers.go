package api

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/shop/api/account"
	"github.com/nbvghost/dandelion/library/shop/api/content"
	"github.com/nbvghost/dandelion/library/shop/api/express"
	"github.com/nbvghost/dandelion/library/shop/api/goods/review"
	"github.com/nbvghost/dandelion/library/shop/api/goods/wish"
	"github.com/nbvghost/dandelion/library/shop/api/index"
	"github.com/nbvghost/dandelion/library/shop/api/journal"
	"github.com/nbvghost/dandelion/library/shop/api/order"
	"github.com/nbvghost/dandelion/library/shop/api/order/cart"
	"github.com/nbvghost/dandelion/library/shop/api/payment/method/paypal"
	"github.com/nbvghost/dandelion/library/shop/api/session"
	"github.com/nbvghost/dandelion/library/shop/api/store"
	"github.com/nbvghost/dandelion/library/shop/api/user"
	"github.com/nbvghost/dandelion/library/shop/api/wx"
)

func Register(route constrain.IRoute) {
	route.RegisterRoute("leave-message", &LeaveMessage{})

	route.RegisterRoute("account/mini_program_login", &account.MiniProgramLogin{})
	route.RegisterRoute("account/qrcode", &account.MiniprogramQRcode{})
	route.RegisterRoute("account/get_login_user", &account.GetLoginUser{})
	route.RegisterRoute("account/get_login_user_phone", &account.GetLoginUserPhone{})
	route.RegisterRoute("account/config", &account.Config{})
	route.RegisterRoute("account/captcha", &account.Captcha{})
	route.RegisterRoute("account/reset-password", &account.ResetPassword{})

	route.RegisterRoute("wx/callback", &wx.Callback{})
	route.RegisterRoute("wx/notify/{OID}", &wx.Notify{})
	route.RegisterRoute("wx/refund-notify/{OID}", &wx.RefundNotify{})
	route.RegisterRoute("wx/message/{OID}", &wx.Message{})
	route.RegisterRoute("wx/token", &wx.Token{})

	route.RegisterRoute("user/level/{UserID}", &user.Level{})
	route.RegisterRoute("user/info", &user.Info{})
	route.RegisterRoute("user/info/DaySign", &user.DaySign{})
	route.RegisterRoute("user/growth/list/{Order}", &user.GrowthList{})
	route.RegisterRoute("user/info/{UserID}", &user.InfoUser{})
	route.RegisterRoute("user/info/sharekey", &user.InfoSharekey{})
	route.RegisterRoute("user/transfers", &user.Transfers{})
	route.RegisterRoute("user/update", &user.Update{})
	route.RegisterRoute("user/upload-avatar", &user.UploadAvatar{})
	route.RegisterRoute("user/user", &user.User{})
	route.RegisterRoute("user/address", &user.Address{})
	route.RegisterRoute("user/review", &user.Review{})
	route.RegisterRoute("user/review-details", &user.ReviewDetails{})

	route.RegisterRoute("goods/wish/goods", &wish.Goods{})
	route.RegisterRoute("goods/review/goods-info", &review.GoodsInfo{}) //来自产品页的评论

	route.RegisterRoute("store/location/list", &store.LocationList{})
	route.RegisterRoute("store/get", &store.Get{})
	route.RegisterRoute("store/get/{StoreID}", &store.GetStore{})
	route.RegisterRoute("store/list/stock", &store.StockList{})
	route.RegisterRoute("store/list/stock/goods/specification/{GoodsID}", &store.StockGoodsList{})
	route.RegisterRoute("store/verification/get/{VerificationNo}", &store.Verification{})
	route.RegisterRoute("store/verification", &store.Verification{})
	route.RegisterRoute("store/supply", &store.Supply{})
	route.RegisterRoute("store/journal/list", &store.JournalList{})
	route.RegisterRoute("store/transfers", &store.Transfers{})
	route.RegisterRoute("store/add/star", &store.AddStar{})

	route.RegisterRoute("content/{ContentItemID}/list/hot", &content.ListHot{})
	route.RegisterRoute("content/{ContentItemID}/list/new", &content.ListNew{})
	route.RegisterRoute("content/article/{ArticleID}", &content.Article{})
	route.RegisterRoute("content/{ContentItemID}/list/subtype", &content.ListContentSubType{})
	route.RegisterRoute("content/{ContentItemID}/related/{ContentSubTypeID}", &content.Related{})
	route.RegisterRoute("content/like", &content.Like{})

	route.RegisterRoute("journal/list/leve", &journal.ListLeve{})
	route.RegisterRoute("journal/list/journal", &journal.ListJournal{})

	route.RegisterRoute("order/buy", &order.Buy{})
	route.RegisterRoute("order/buy/collage", &order.BuyCollage{})
	{
		route.RegisterRoute("order/cart/delete", &cart.Delete{})
		route.RegisterRoute("order/cart/change", &cart.Change{})
		route.RegisterRoute("order/cart/list", &cart.List{})
		route.RegisterRoute("order/cart/add", &cart.Add{})
	}
	route.RegisterRoute("order/confirm/list", &order.ConfirmList{})
	route.RegisterRoute("order/create-orders", &order.CreateOrders{})
	route.RegisterRoute("order/info-orders", &order.InfoOrders{})
	route.RegisterRoute("order/wxpay/package", &order.WXPayPackage{})
	route.RegisterRoute("order/wxpay/alone", &order.WXPayAlone{})
	route.RegisterRoute("order/list", &order.List{})
	route.RegisterRoute("order/{ID:[0-9]+}/get", &order.GetOrder{})
	route.RegisterRoute("order/collage/record", &order.CollageRecord{})
	route.RegisterRoute("order/change", &order.Change{})
	route.RegisterRoute("order/express/info", &order.ExpressInfo{})

	route.RegisterRoute("express/delivery-list", &express.DeliveryList{})

	route.RegisterRoute("index/goods_type/list", &index.GoodsTypeList{})
	route.RegisterRoute("index/goods_type/child/{GoodsTypeID}/list", &index.GoodsTypeChildList{})
	route.RegisterRoute("index/goods/child/{GoodsTypeID}/{GoodsTypeChildID}/list", &index.GoodsChild{})
	route.RegisterRoute("index/goods/get/{ID}", &index.Goods{})
	route.RegisterRoute("index/goods/hot/list", &index.GoodsHotList{})
	route.RegisterRoute("index/goods/trending/list", &index.GoodsTrendingList{})
	route.RegisterRoute("index/goods/all/list", &index.GoodsAllList{})
	route.RegisterRoute("index/score_goods/list", &index.ScoreGoodsList{})
	route.RegisterRoute("index/score_goods/exchange/{ScoreGoodsID}", &index.ScoreGoodsExchange{})
	route.RegisterRoute("index/share/score", &index.ShareScore{})
	route.RegisterRoute("index/card/list", &index.CardList{})
	route.RegisterRoute("index/card/get/{CardItemID}", &index.CardGet{})
	route.RegisterRoute("index/verification/get/{VerificationNo}", &index.VerificationGet{})
	route.RegisterRoute("index/read/share/key", &index.ReadShareKey{})
	route.RegisterRoute("index/configuration/list", &index.ConfigurationList{})
	route.RegisterRoute("index/push-data", &index.PushData{})

	//session
	route.RegisterRoute("session/index", &session.Index{})

	//payment
	route.RegisterRoute("payment/method/paypal/checkout-orders", &paypal.CheckoutOrders{})
	route.RegisterRoute("payment/method/paypal/capture/{PaypalOrderID}", &paypal.Capture{})

}
