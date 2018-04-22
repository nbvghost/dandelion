package district

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nbvghost/gweb"

	"dandelion/app/service"
)

type Controller struct {
	gweb.BaseController
}

func (c *Controller) Apply() {
	c.AddHandler(gweb.ALLMethod("get", getAction))

}
func getAction(context *gweb.Context) gweb.Result {

	code := context.Request.URL.Query().Get("code")
	codes := strings.Split(code, ",")
	fmt.Println(codes)
	len := len(codes)
	if strings.EqualFold(code, "") {
		len = 0
	}
	switch len {
	case 0:
		return &gweb.JsonResult{Data: service.District.ProvinceDao.ListP(service.Orm)}
	case 1:
		p, _ := strconv.Atoi(codes[0])
		return &gweb.JsonResult{Data: service.District.CityDao.ListC(service.Orm, p)}
	case 2:
		p, _ := strconv.Atoi(codes[0])
		c, _ := strconv.Atoi(codes[1])
		return &gweb.JsonResult{Data: service.District.AreaDao.ListA(service.Orm, p, c)}

	}
	return &gweb.JsonResult{}
}
