package api

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/admin/api/account"
	"github.com/nbvghost/dandelion/library/admin/api/activity/collage"
	"github.com/nbvghost/dandelion/library/admin/api/activity/fullcut"
	"github.com/nbvghost/dandelion/library/admin/api/activity/score_goods"
	"github.com/nbvghost/dandelion/library/admin/api/activity/timesell"
	"github.com/nbvghost/dandelion/library/admin/api/activity/voucher"
	"github.com/nbvghost/dandelion/library/admin/api/admin"
	"github.com/nbvghost/dandelion/library/admin/api/carditem"
	"github.com/nbvghost/dandelion/library/admin/api/company"
	"github.com/nbvghost/dandelion/library/admin/api/configuration"
	"github.com/nbvghost/dandelion/library/admin/api/content"
	"github.com/nbvghost/dandelion/library/admin/api/content/content_item"
	"github.com/nbvghost/dandelion/library/admin/api/content/content_sub_type"
	"github.com/nbvghost/dandelion/library/admin/api/content/customize/field"
	"github.com/nbvghost/dandelion/library/admin/api/content/leave_message"
	"github.com/nbvghost/dandelion/library/admin/api/express"
	"github.com/nbvghost/dandelion/library/admin/api/express/template"
	"github.com/nbvghost/dandelion/library/admin/api/file"
	"github.com/nbvghost/dandelion/library/admin/api/goods"
	"github.com/nbvghost/dandelion/library/admin/api/notify"
	"github.com/nbvghost/dandelion/library/admin/api/order"
	"github.com/nbvghost/dandelion/library/admin/api/order/shipping"
	"github.com/nbvghost/dandelion/library/admin/api/store"
	"github.com/nbvghost/dandelion/library/admin/api/store/store_stock"
	"github.com/nbvghost/dandelion/library/admin/api/store_journal"
)

func Register(route constrain.IRoute) {

	//adminController := &admin.Controller{}
	//adminController.Interceptors.Set(&admin.Interceptor{})
	//adminController := gweb.NewController("admin", "")
	//adminController.NewController("template").DefaultHandle(&admin.Index{})
	//adminController.AddInterceptor(&admin.Interceptor{})

	route.RegisterRoute("account/login", &account.Login{})
	route.RegisterRoute("account/register", &account.Register{})

	route.RegisterRoute("notify/login", &notify.Login{})

	route.RegisterRoute("heartbeat", &Heartbeat{})
	route.RegisterRoute("translate", &Translate{})

	route.RegisterRoute("goods/attributes", &goods.Attributes{})
	route.RegisterRoute("goods/attributes-group", &goods.AttributesGroup{})
	route.RegisterRoute("goods/change-goods", &goods.ChangeGoods{})
	route.RegisterRoute("goods/specification", &goods.Specification{})
	route.RegisterRoute("goods/add-goods", &goods.AddGoods{})
	route.RegisterRoute("goods/list-goods-type", &goods.ListGoodsType{})
	route.RegisterRoute("goods/list-goods-type-child-id", &goods.ListGoodsTypeChildID{})
	route.RegisterRoute("goods/add-goods-type-child", &goods.AddGoodsTypeChild{})
	route.RegisterRoute("goods/add-goods-type", &goods.AddGoodsType{})
	route.RegisterRoute("goods/change-goods-type", &goods.ChangeGoodsType{})
	route.RegisterRoute("goods/list-goods-type-child", &goods.ListGoodsTypeChild{})
	route.RegisterRoute("goods/del-goods-type", &goods.DelGoodsType{})
	route.RegisterRoute("goods/list-goods-type-all", &goods.ListGoodsTypeAll{})
	route.RegisterRoute("goods/activity-goods", &goods.ActivityGoods{})
	route.RegisterRoute("goods/list-goods", &goods.ListGoods{})
	route.RegisterRoute("goods/change-goods-type-child", &goods.ChangeGoodsTypeChild{})
	route.RegisterRoute("goods/get-goods", &goods.GetGoods{})
	route.RegisterRoute("goods/list-specification", &goods.ListSpecification{})
	route.RegisterRoute("goods/get-goods-type", &goods.GetGoodsType{})
	route.RegisterRoute("goods/get-goods-type-child", &goods.GetGoodsTypeChild{})
	route.RegisterRoute("goods/del-goods-type-child", &goods.DelGoodsTypeChild{})
	route.RegisterRoute("goods/del-goods", &goods.DelGoods{})
	route.RegisterRoute("goods/sku-label", &goods.SkuLabel{})
	route.RegisterRoute("goods/sku-label-data", &goods.SkuLabelData{})

	route.RegisterRoute("order/list", &order.List{})
	route.RegisterRoute("order/change", &order.Change{})
	route.RegisterRoute("order/query", &order.Query{})
	route.RegisterRoute("order/shipping/shipping", &shipping.Shipping{})

	route.RegisterRoute("store_journal/list", &store_journal.List{})
	route.RegisterRoute("configuration/list", &configuration.List{})
	route.RegisterRoute("configuration/change", &configuration.Change{})
	route.RegisterRoute("carditem/list", &carditem.List{})

	route.RegisterRoute("admin/situation", &admin.Situation{})
	route.RegisterRoute("admin/index", &admin.Admin{})
	route.RegisterRoute("admin/list", &admin.List{})
	route.RegisterRoute("admin/login", &admin.Login{})
	route.RegisterRoute("admin/login-out", &admin.LoginOut{})
	route.RegisterRoute("admin/role-list", &admin.RoleList{})
	route.RegisterRoute("admin/authority/{ID}", &admin.AuthorityID{})
	route.RegisterRoute("admin/{ID}", &admin.ID{})

	route.RegisterRoute("express/template/save", &template.Save{})
	route.RegisterRoute("express/template/list", &template.List{})
	route.RegisterRoute("express/template/{ID:[0-9]+}", &template.ID{})
	route.RegisterRoute("express/template/table-list", &template.TableList{})
	route.RegisterRoute("express/delivery-list", &express.DeliveryList{})

	route.RegisterRoute("table-list", &content.TableList{})
	route.RegisterRoute("save", &content.Save{})
	route.RegisterRoute("get", &content.Get{})
	route.RegisterRoute("single/get/{ContentItemID}/{ContentSubTypeID}", &content.SingleGetContentItemIDContentSubTypeID{})

	route.RegisterRoute("content/delete", &content.Delete{})
	route.RegisterRoute("content/multi-get/{ID}", &content.MultiGetID{})
	route.RegisterRoute("content/table-list", &content.TableList{})
	route.RegisterRoute("content/type-list", &content.TypeList{})
	route.RegisterRoute("content/config", &content.Config{})
	route.RegisterRoute("content/sub-type", &content.SubType{})
	route.RegisterRoute("content/save", &content.Save{})
	route.RegisterRoute("content/list", &content.List{})
	route.RegisterRoute("content/change", &content.Change{})
	route.RegisterRoute("content/cache", &content.Cache{})
	route.RegisterRoute("content/single-get/{ContentItemID}/{ContentSubTypeID}", &content.SingleGetContentItemIDContentSubTypeID{}) //single-get-ContentItemID-ContentSubTypeID.go
	route.RegisterRoute("content/content_item/add", &content_item.Add{})
	route.RegisterRoute("content/content_item/{ContentItemID:[0-9]+}", &content_item.ContentItemID{})
	route.RegisterRoute("content/content_item/list", &content_item.List{})
	route.RegisterRoute("content/content_item/index/{ContentItemID}", &content_item.IndexContentItemID{})
	route.RegisterRoute("content/content_item/hide/{ContentItemID}", &content_item.HideContentItemID{})
	route.RegisterRoute("content/content_item/show-at-home", &content_item.ShowAtHome{})
	route.RegisterRoute("content/content_item/config", &content_item.Config{})

	route.RegisterRoute("content/leave_message/list", &leave_message.List{})

	route.RegisterRoute("content/content-sub-type/list/{ContentItemID}", &content_sub_type.ListContentItemID{})
	route.RegisterRoute("content/content-sub-type/child-list/{ContentItemID}/{ParentContentSubTypeID}", &content_sub_type.ChildListContentItemIDParentContentSubTypeID{})
	route.RegisterRoute("content/content-sub-type/all-list/{ContentItemID}", &content_sub_type.AllListContentItemID{})
	route.RegisterRoute("content/content-sub-type/get/{ContentSubTypeID}", &content_sub_type.GetContentSubTypeID{})
	route.RegisterRoute("content/content-sub-type/{ID}", &content_sub_type.ID{})

	route.RegisterRoute("content/customize/field/group", &field.Group{})
	route.RegisterRoute("content/customize/field/field", &field.Field{})
	route.RegisterRoute("content/customize/field/sort-field", &field.SortField{})
	route.RegisterRoute("content/customize/field/list-field", &field.ListField{})
	route.RegisterRoute("content/customize/field/sync", &field.Sync{})

	route.RegisterRoute("company/info", &company.Info{})
	route.RegisterRoute("store/add", &store.Add{})
	route.RegisterRoute("store/list", &store.List{})
	route.RegisterRoute("store/stock", &store.Stock{})
	route.RegisterRoute("store/{ID}", &store.ID{})

	route.RegisterRoute("store/store_stock/exist-goods/{StoreID}", &store_stock.ExistGoodsStoreID{})
	route.RegisterRoute("store/store_stock/list/{StoreID}/{GoodsID}", &store_stock.ListStoreIDGoodsID{})
	route.RegisterRoute("store/store_stock/list", &store_stock.List{})
	route.RegisterRoute("store/store_stock/{ID}", &store_stock.ID{})

	route.RegisterRoute("activity/collage/save", &collage.Save{})
	route.RegisterRoute("activity/collage/change", &collage.Change{})
	route.RegisterRoute("activity/collage/{Hash}", &collage.Hash{})
	route.RegisterRoute("activity/collage/list", &collage.List{})
	route.RegisterRoute("activity/collage/goods-{Hash}-list", &collage.GoodsHashList{})
	route.RegisterRoute("activity/collage/goods-{GoodsID}", &collage.GoodsGoodsID{})
	route.RegisterRoute("activity/collage/goods-add", &collage.GoodsAdd{})
	route.RegisterRoute("activity/collage/{ID}", &collage.ID{})
	route.RegisterRoute("activity/fullcut/save", &fullcut.Save{})
	route.RegisterRoute("activity/fullcut/{ID}", &fullcut.ID{})
	route.RegisterRoute("activity/fullcut/list", &fullcut.List{})
	route.RegisterRoute("activity/score_goods/index", &score_goods.Index{})
	route.RegisterRoute("activity/score_goods/list", &score_goods.List{})
	route.RegisterRoute("activity/score_goods/{ID}", &score_goods.ID{})
	route.RegisterRoute("activity/timesell/save", &timesell.Save{})
	route.RegisterRoute("activity/timesell/change", &timesell.Change{})
	route.RegisterRoute("activity/timesell/{Hash}", &timesell.Hash{})
	route.RegisterRoute("activity/timesell/list", &timesell.List{})
	route.RegisterRoute("activity/timesell/goods-{Hash}-list", &timesell.GoodsHashList{})
	route.RegisterRoute("activity/timesell/goods-{GoodsID}", &timesell.GoodsGoodsID{})
	route.RegisterRoute("activity/timesell/goods-add", &timesell.GoodsAdd{})
	route.RegisterRoute("activity/timesell/{ID}", &timesell.ID{})
	route.RegisterRoute("activity/voucher/index", &voucher.Index{})
	route.RegisterRoute("activity/voucher/list", &voucher.List{})
	route.RegisterRoute("activity/voucher/{ID}", &voucher.ID{})

	route.RegisterRoute("file/up", &file.Up{})
	route.RegisterRoute("file/load", &file.Load{})
	route.RegisterRoute("file/ckeditor-up", &file.CKEditorUp{})

}
