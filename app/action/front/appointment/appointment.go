package appointment

import (
	"github.com/nbvghost/dandelion/app/result"
	"net/http"
	"net/url"

	"github.com/nbvghost/gweb"

	"github.com/nbvghost/dandelion/app/play"
)

type InterceptorFile struct {
}

func (this InterceptorFile) Execute(context *gweb.Context) bool {

	//util.Trace(context.Session,"context.Session")
	if context.Session.Attributes.Get(play.SessionManager) != nil || context.Session.Attributes.Get(play.SessionAdmin) != nil {

		return true
	} else {
		redirect := ""
		if len(context.Request.URL.Query().Encode()) == 0 {
			redirect = context.Request.URL.Path
		} else {
			redirect = context.Request.URL.Path + "?" + context.Request.URL.Query().Encode()
		}
		//fmt.Println(url.QueryEscape(redirect))
		http.Redirect(context.Response, context.Request, "/account/loginManagerPage?redirect="+url.QueryEscape(redirect), http.StatusFound)
		return false

	}
}

type Controller struct {
	gweb.BaseController
}

func (sub *Controller) Init() {
	//i.Interceptors.Add(&InterceptorFile{})
	//sub.AddHandler(gweb.ALLMethod(":id/index", indexPage))
	sub.AddHandler(gweb.ALLMethod(":id/action/add", actionAddAction))
	//sub.AddHandler(gweb.ALLMethod("/", loadFilePage))
}
func actionAddAction(context *gweb.Context) gweb.Result {
	/* ID, _ := strconv.ParseUint(context.PathParams["id"], 10, 64)
	appointment := service.Appointment.GetAppointment(ID)
	context.Request.ParseForm()

	Count, _ := strconv.ParseUint(context.Request.Form.Get("Count"), 10, 64)

	var carts map[uint64]*dao.TempOrderPack

	if context.Session.Attributes.Get(play.SessionCart) != nil {
		carts = context.Session.Attributes.Get(play.SessionCart).(map[uint64]*dao.TempOrderPack)
	} else {
		carts = make(map[uint64]*dao.TempOrderPack)
	}

	if carts[appointment.CompanyID] == nil {
		shop := service.Company.GetCompany(appointment.CompanyID)
		carts[appointment.CompanyID] = &dao.TempOrderPack{appointment.CompanyID, shop.Name, make([]*dao.Order, 0)}
	}

	tempOrderPack := carts[appointment.CompanyID]

	have := false
	for i := 0; i < len(tempOrderPack.Orders); i++ {
		_order := tempOrderPack.Orders[i]
		if _order.OID == appointment.ID {
			have = true
			_order.Count = _order.Count + 1
			_order.Total = _order.Price * _order.Count
		}

	}
	if have == false {

		order := &dao.Order{}
		order.Count = Count
		order.OID = appointment.ID
		order.Title = appointment.Name
		order.Type = 1
		order.CompanyID = appointment.CompanyID
		//content, _ := json.Marshal(appointment)
		order.CacheContent = appointment
		order.UserID = 4
		order.Price = appointment.Price
		order.Total = order.Price * Count

		tempOrderPack.Orders = append(tempOrderPack.Orders, order)
	}
	context.Session.Attributes.Put(play.SessionCart, carts) */
	return &gweb.JsonResult{Data: (&result.ActionResult{}).Smart(result.ActionOK, "添加成功", "添加失败")}
}

/* func indexPage(context *gweb.Context) gweb.Result {
	 id, _ := strconv.ParseUint(context.PathParams["id"], 10, 64)
	appointment := service.Appointment.GetAppointment(id)
	shop := service.Company.GetCompany(appointment.CompanyID)

	return &gweb.HTMLResult{Name: "appointment/index", Params: map[string]interface{}{"appointment": appointment, "shop": shop}}
}
func loadFilePage(context *gweb.Context) gweb.Result {

	return service.File.LoadAction(context)
}
*/
