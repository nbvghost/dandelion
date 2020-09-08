package order

import (
	"fmt"

	"github.com/nbvghost/gweb"

	"github.com/nbvghost/dandelion/app/util"
	"github.com/nbvghost/dandelion/app/util/wxpay"
)

type Controller struct {
	gweb.BaseController
}

func (i *Controller) Init() {
	//Index.RequestMapping = make(map[string]mvc.Function)
	i.AddHandler(gweb.ALLMethod("", defaultPage))
	i.AddHandler(gweb.ALLMethod("orderAction", orderAction))
	i.AddHandler(gweb.ALLMethod("index", indexPage))
	i.AddHandler(gweb.ALLMethod("pay", payPage))

}
func orderAction(context *gweb.Context) gweb.Result {
	action := context.Request.URL.Query().Get("action")
	result := &gweb.JsonResult{}
	switch action {
	case "list":
		//var carts map[uint64]*dao.TempOrderPack
		/* carts := context.Session.Attributes.Get(play.SessionCart)
		if carts == nil {
			carts = make(map[uint64]*dao.TempOrderPack)
		} */

		/*_carts := carts.(map[uint64]*dao.TempOrderPack)
		for _, value := range _carts {
			fmt.Println(value)
		}*/
		//result.Data = (&result.ActionResult{}).SmartSuccessData(carts)
	case "Pay":
		/* ShopID, _ := strconv.ParseUint(context.Request.URL.Query().Get("ShopID"), 10, 64)

		carts := context.Session.Attributes.Get(play.SessionCart)
		if carts == nil {
			carts = make(map[uint64]*dao.TempOrderPack)
		}
		cart := carts.(map[uint64]*dao.TempOrderPack)[ShopID]

		Position := context.Request.URL.Query().Get("Position")
		Tip := context.Request.URL.Query().Get("Tip")

		orderPack := &dao.OrderPack{}
		orderPack.UserID = 5
		orderPack.Status = 0
		orderPack.OrderNo = tool.UUID()
		orderPack.CompanyID = 5
		orderPack.CreatedAt = time.Now()
		//orderPack.ClosedAt
		//orderPack.PayAt
		orderPack.Position = Position
		orderPack.Tip = Tip

		b, err := json.Marshal(cart.Orders)
		glog.Error(err)
		orderPack.OrderList = string(b)
		service.OrderPack.Add(service.Orm, orderPack)

		goodsDetail := make([]wxpay.WXGoodsDetail, 0)
		for _, value := range cart.Orders {
			goodsDetail = append(goodsDetail, wxpay.WXGoodsDetail{
				strconv.FormatUint(value.OID, 10),
				value.Title,
				strconv.FormatUint(value.Count, 10),
				strconv.FormatUint(value.Price, 10),
			})
			value.Total = value.Price * value.Count
			orderPack.Total = orderPack.Total + value.Total
		}
		detail := wxpay.WXDetail{goodsDetail}
		b, err = json.Marshal(detail)

		wxpay.OrderJS(orderPack.OrderNo, "困国国国", string(b), util.GetHost(context), "", "192.168.1.122", orderPack.Total) */

	case "Count":
		/* index, _ := strconv.Atoi(context.Request.URL.Query().Get("index"))
		value, _ := strconv.Atoi(context.Request.URL.Query().Get("value"))
		ShopID, _ := strconv.ParseUint(context.Request.URL.Query().Get("ShopID"), 10, 64)

		carts := context.Session.Attributes.Get(play.SessionCart)
		if carts == nil {
			carts = make(map[uint64]*dao.TempOrderPack)
		}
		cart := carts.(map[uint64]*dao.TempOrderPack)[ShopID]

		if index < len(cart.Orders) {
			if value < 0 {
				vp := cart.Orders[index].Count - 1
				if cart.Orders[index].Count == 1 {
					vp = 1
				}
				cart.Orders[index].Count = vp
			} else {
				vp := cart.Orders[index].Count + 1
				cart.Orders[index].Count = vp
			}
			//context.Session.Attributes.Put(play.SessionCart, _carts)
		} */

	case "del":
		/* index, _ := strconv.Atoi(context.Request.URL.Query().Get("index"))
		ShopID, _ := strconv.ParseUint(context.Request.URL.Query().Get("ShopID"), 10, 64)
		carts := context.Session.Attributes.Get(play.SessionCart)
		if carts == nil {
			carts = make(map[uint64]*dao.TempOrderPack)
		}
		cart := carts.(map[uint64]*dao.TempOrderPack)[ShopID]
		if index < len(cart.Orders) {
			//_carts := carts.(map[uint64]*dao.TempOrderPack)
			//s := []int{11, 22, 33, 44, 55, 66} // 原始slice
			//i := 2                        // 待删除项索引
			cart.Orders = append(cart.Orders[:index], cart.Orders[index+1:]...) // 最后面的“...”不能省略
			//context.Session.Attributes.Put(play.SessionCart, _carts)
		}
		//OrderNo := context.Request.URL.Query().Get("OrderNo") */
	}

	return result
}
func defaultPage(context *gweb.Context) gweb.Result {

	return &gweb.RedirectToUrlResult{"index"}
}
func payPage(context *gweb.Context) gweb.Result {

	d := make(map[string]interface{})

	appId, timestamp, nonceStr, signature := wxpay.GetWXJSConfig(util.GetFullUrl(context.Request))
	d["appId"] = appId
	d["timestamp"] = timestamp
	d["nonceStr"] = nonceStr
	d["signature"] = signature

	fmt.Println(d)
	fmt.Println(appId)
	fmt.Println(timestamp)
	fmt.Println(nonceStr)
	fmt.Println(signature)

	return &gweb.HTMLResult{Params: d}
}
func indexPage(context *gweb.Context) gweb.Result {

	return &gweb.RedirectToUrlResult{"pay"}
}
