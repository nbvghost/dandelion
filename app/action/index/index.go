package index

import "github.com/nbvghost/gweb"

type Controller struct {
	gweb.BaseController
}

func (c *Controller) Apply() {
	//Index.RequestMapping = make(map[string]mvc.Function)
	c.AddHandler(gweb.ALLMethod("", sdfsd))
	c.AddHandler(gweb.ALLMethod("*", indexPage))
	c.AddHandler(gweb.ALLMethod("index", indexPage))
	c.AddHandler(gweb.ALLMethod("mPDhoTorHe.txt", dsfdsfdsPage))
	c.AddHandler(gweb.ALLMethod("P962URPfYr.txt", dsfdsfsdfsddsPage))

}
func sdfsd(context *gweb.Context) gweb.Result {

	return &gweb.RedirectToUrlResult{"index"}
}

//6c0420c5e926a2ac8d56aa4192ab10fa
func indexPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}
func dsfdsfdsPage(context *gweb.Context) gweb.Result {

	return &gweb.TextResult{Data: "6c0420c5e926a2ac8d56aa4192ab10fa"}
}
func dsfdsfsdfsddsPage(context *gweb.Context) gweb.Result {

	return &gweb.TextResult{Data: "15c9e3b7158aeefc0444f9fe862954b3"}
}
