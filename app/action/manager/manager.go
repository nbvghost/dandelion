package manager

import (
	"net/url"

	"github.com/nbvghost/glog"

	"github.com/nbvghost/dandelion/app/util"

	"encoding/json"
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/service"
	"github.com/nbvghost/dandelion/app/service/dao"
	"strconv"

	"github.com/nbvghost/gweb"
)

type InterceptorManager struct {
}

//Execute(Session *Session,Request *http.Request)(bool,Result)
func (this InterceptorManager) Execute(context *gweb.Context) (bool, gweb.Result) {

	//util.Trace(context.Session,"context.Session")
	if context.Session.Attributes.Get(play.SessionManager) == nil {
		//http.SetCookie(context.Response, &http.Cookie{Name: "UID", MaxAge:-1, Path: "/"})
		//fmt.Println(context.Request.URL.Path)
		//fmt.Println(context.Request.URL.Query().Encode())
		redirect := ""
		if len(context.Request.URL.Query().Encode()) == 0 {
			redirect = context.Request.URL.Path
		} else {
			redirect = context.Request.URL.Path + "?" + context.Request.URL.Query().Encode()
		}
		//fmt.Println(url.QueryEscape(redirect))
		//http.Redirect(context.Response, context.Request, "/account/loginManagerPage?redirect="+url.QueryEscape(redirect), http.StatusFound)
		return false, &gweb.RedirectToUrlResult{Url: "/account/loginManagerPage?redirect=" + url.QueryEscape(redirect)}
	} else {
		return true, nil
	}
}

type Controller struct {
	gweb.BaseController
	Article       service.ArticleService
	Admin         service.AdminService
	User          service.UserService
	Rank          service.RankService
	Configuration service.ConfigurationService
	GiveVoucher   service.GiveVoucherService
	Voucher       service.VoucherService
	Goods         service.GoodsService
}

func (controller *Controller) Apply() {
	//Index.RequestMapping = make(map[string]mvc.Function)
	controller.Interceptors.Add(&InterceptorManager{})

	controller.AddHandler(gweb.ALLMethod("index", controller.indexPage))
	controller.AddHandler(gweb.ALLMethod("articlePage", controller.articlePage))
	controller.AddHandler(gweb.ALLMethod("article", controller.articleAction))
	controller.AddHandler(gweb.ALLMethod("admin", controller.adminAction))
	controller.AddHandler(gweb.ALLMethod("add_article", controller.addArticlePage))

	controller.AddHandler(gweb.POSMethod("rank/add", controller.rankAddAction))
	controller.AddHandler(gweb.POSMethod("rank/list", controller.rankListAction))
	controller.AddHandler(gweb.DELMethod("rank/{RankID}", controller.rankDeleteAction))

	controller.AddHandler(gweb.POSMethod("user/all/list", controller.ListAllTableDatas))

	controller.AddHandler(gweb.POSMethod("configuration/list", controller.configurationListAction))
	controller.AddHandler(gweb.POSMethod("configuration/change", controller.configurationChangeAction))
	//this.AddHandler(gweb.ALLMethod("categoryAction", this.categoryAction))
	controller.AddHandler(gweb.POSMethod("give_voucher/save", controller.giveVoucherSaveAction))
	controller.AddHandler(gweb.POSMethod("give_voucher/list", controller.giveVoucherListAction))

	controller.AddHandler(gweb.ALLMethod("goods", controller.GoodsAction))

}
func (controller *Controller) GoodsAction(context *gweb.Context) gweb.Result {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	action := context.Request.URL.Query().Get("action")
	Orm := dao.Orm()
	switch action {

	case "del_goods_type":
		ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("ID"), 10, 64)
		return &gweb.JsonResult{Data: controller.Goods.DeleteGoodsType(ID)}
	case "add_goods_type":
		item := &dao.GoodsType{}
		//item.OID = company.ID
		err := util.RequestBodyToJSON(context.Request.Body, item)
		glog.Trace(err)

		//fmt.Println(item)
		err = controller.Goods.Add(Orm, item)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "添加成功", nil)}
	case "change_goods_type":
		item := &dao.GoodsType{}
		err := util.RequestBodyToJSON(context.Request.Body, item)
		glog.Trace(err)
		err = controller.Goods.ChangeModel(Orm, item.ID, &dao.GoodsType{Name: item.Name})
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}

	case "get_goods_type_child":
		ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("ID"), 10, 64)
		var goods dao.GoodsTypeChild
		controller.Goods.Get(Orm, ID, &goods)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", goods)}
	case "del_goods_type_child":
		ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("ID"), 10, 64)
		return &gweb.JsonResult{Data: controller.Goods.DeleteGoodsTypeChild(ID)}
	case "add_goods_type_child":
		item := &dao.GoodsTypeChild{}
		err := util.RequestBodyToJSON(context.Request.Body, item)
		glog.Trace(err)
		//fmt.Println(item)
		err = controller.Goods.Add(Orm, item)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "添加成功", nil)}
	case "change_goods_type_child":
		item := &dao.GoodsTypeChild{}
		err := util.RequestBodyToJSON(context.Request.Body, item)
		glog.Trace(err)
		err = controller.Goods.ChangeModel(Orm, item.ID, &dao.GoodsTypeChild{Name: item.Name, Image: item.Image})
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
	case "list_goods_type_child":
		var gts []dao.GoodsTypeChild
		controller.Goods.FindAll(Orm, &gts)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", gts)}
	case "list_goods_type_child_id":
		ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("ID"), 10, 64)
		gts := controller.Goods.ListAllGoodsTypeChild(ID)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", gts)}

	case "list_goods_type_all":
		//gts := controller.Goods.ListGoodsType()
		//return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", gts)}
	case "list_goods_type":
		dts := &dao.Datatables{}
		util.RequestBodyToJSON(context.Request.Body, dts)
		draw, recordsTotal, recordsFiltered, list := controller.Goods.DatatablesListOrder(Orm, dts, &[]dao.GoodsType{}, 0, "")
		return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}

	}

	return &gweb.JsonResult{}
}
func (controller *Controller) giveVoucherSaveAction(context *gweb.Context) gweb.Result {

	item := dao.GiveVoucher{}
	util.RequestBodyToJSON(context.Request.Body, &item)
	err := controller.GiveVoucher.SaveItem(item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "添加成功", nil)}
}
func (controller *Controller) giveVoucherListAction(context *gweb.Context) gweb.Result {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := controller.GiveVoucher.DatatablesListOrder(Orm, dts, &[]dao.GiveVoucher{}, 0, "")
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}
func (controller *Controller) configurationChangeAction(context *gweb.Context) gweb.Result {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	item := dao.Configuration{}
	util.RequestBodyToJSON(context.Request.Body, &item)
	err := controller.Configuration.ChangeConfiguration(0, item.K, item.V)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
}
func (controller *Controller) configurationListAction(context *gweb.Context) gweb.Result {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	var ks []uint64
	util.RequestBodyToJSON(context.Request.Body, &ks)
	list := controller.Configuration.GetConfigurations(0, ks)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", list)}
}
func (controller *Controller) rankAddAction(context *gweb.Context) gweb.Result {

	rank := dao.Rank{}
	util.RequestBodyToJSON(context.Request.Body, &rank)
	err := controller.Rank.AddRank(rank)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "添加成功", nil)}
}
func (controller *Controller) rankListAction(context *gweb.Context) gweb.Result {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := controller.Rank.DatatablesListOrder(Orm, dts, &[]dao.Rank{}, 0, "")
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}
func (controller *Controller) rankDeleteAction(context *gweb.Context) gweb.Result {

	RankID, _ := strconv.ParseUint(context.PathParams["RankID"], 10, 64)

	err := controller.Rank.Delete(dao.Orm(), &dao.Rank{}, RankID)

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}
func (controller *Controller) ListAllTableDatas(context *gweb.Context) gweb.Result {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := controller.User.DatatablesListOrder(Orm, dts, &[]dao.User{}, 0, "")
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}
func (controller *Controller) addArticlePage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{Params: nil}
}

/*func (this *Controller) categoryAction(context *gweb.Context) gweb.Result {
	action := context.Request.URL.Query().Get("action")
	result := &gweb.JsonResult{}
	switch action {
	case "list":
		result.Data = &dao.ActionStatus{true, "ok", this.ArticleType.FindCategory()}
	case "add":
		context.Request.ParseForm()
		label := context.Request.Form.Get("label")
		_, su := this.ArticleType.AddCategory(label)
		result.Data = (&dao.ActionStatus{}).Smart(su, "添加成功", "添加失败,")
	case "del":
		context.Request.ParseForm()
		id, _ := strconv.ParseUint(context.Request.Form.Get("id"), 10, 64)
		this.ArticleType.DelCategory(id)
		result.Data = &dao.ActionStatus{true, "删除成功", nil}
	}

	return result
}*/

func (controller *Controller) adminAction(context *gweb.Context) gweb.Result {
	action := context.Request.URL.Query().Get("action")
	result := &gweb.JsonResult{}
	switch action {
	//list del
	case play.ActionKey_list:
		result.Data = &dao.ActionStatus{true, "ok", controller.Admin.FindAdmin()}
	case play.ActionKey_del:
		ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("id"), 10, 64)
		result.Data = controller.Admin.DelAdmin(ID)
	case play.ActionKey_get:
		ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("id"), 10, 64)
		admin := controller.Admin.GetAdmin(ID)
		result.Data = &dao.ActionStatus{true, "ok", admin}
	case play.ActionKey_change:
		context.Request.ParseForm()
		Name := context.Request.Form.Get("Account")
		Password := context.Request.Form.Get("PassWord")
		//Email := context.Request.Form.Get("Email")
		//Tel := context.Request.Form.Get("Tel")
		ID, _ := strconv.ParseUint(context.Request.Form.Get("ID"), 10, 64)

		if err := controller.Admin.ChangeAdmin(Name, Password, ID); err != nil {
			result.Data = &dao.ActionStatus{false, "修改失败", nil}
		} else {
			result.Data = &dao.ActionStatus{true, "修改成功", nil}
		}

	case play.ActionKey_add:
		context.Request.ParseForm()
		Name := context.Request.Form.Get("Account")
		Password := context.Request.Form.Get("PassWord")
		Domain := context.Request.Form.Get("Domain")
		//Email := context.Request.Form.Get("Email")
		//Tel := context.Request.Form.Get("Tel")
		result.Data = controller.Admin.AddAdmin(Name, Password, Domain)
	}
	return result
}
func (controller *Controller) articleAction(context *gweb.Context) gweb.Result {

	action := context.Request.URL.Query().Get("action")
	result := &gweb.JsonResult{}
	switch action {
	case "listByCategory":
		ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("id"), 10, 64)
		d := controller.Article.FindArticleByContentSubTypeID(ID)
		result.Data = &dao.ActionStatus{Success: true, Message: "SUCCESS", Data: d}
	case "add":
		context.Request.ParseForm()
		jsonText := context.Request.Form.Get("json")
		article := &dao.Article{}
		err := json.Unmarshal([]byte(jsonText), article)
		glog.Error(err)
		result.Data = controller.Article.AddArticle(article)
	case "change":
		context.Request.ParseForm()
		jsonText := context.Request.Form.Get("json")
		article := &dao.Article{}
		err := json.Unmarshal([]byte(jsonText), article)
		glog.Error(err)
		controller.Article.ChangeArticle(article)
		result.Data = &dao.ActionStatus{true, "SUCCESS", article}
	case "one":
		ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("id"), 10, 64)
		article := controller.Article.GetArticle(ID)
		result.Data = &dao.ActionStatus{true, "SUCCESS", article}
	case "del":
		context.Request.ParseForm()
		id, _ := strconv.ParseUint(context.Request.Form.Get("id"), 10, 64)
		controller.Article.DelArticle(id)
		result.Data = &dao.ActionStatus{true, "删除成功", nil}
	}

	return result
}
func (controller *Controller) articlePage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}

func (controller *Controller) indexPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}

func (controller *Controller) defaultPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}
func (controller *Controller) rootPage(context *gweb.Context) gweb.Result {

	return &gweb.RedirectToUrlResult{"index"}
}
