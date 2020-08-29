package district

import (
	"github.com/nbvghost/dandelion/app/service"

	"strings"

	"github.com/nbvghost/gweb"
)

type Controller struct {
	gweb.BaseController
	DistrictService service.DistrictService
}

func (controller *Controller) Init() {
	controller.AddHandler(gweb.GETMethod("get", controller.getAction))

}
func (controller *Controller) getAction(context *gweb.Context) gweb.Result {

	code := context.Request.URL.Query().Get("code")
	codes := strings.Split(code, ",")
	//fmt.Println(codes)
	len := len(codes)
	if strings.EqualFold(code, "") {
		len = 0
	}
	switch len {
	case 0:
		//return &gweb.JsonResult{Data: controller.DistrictService.ProvinceDao.ListP(service.Orm)}
	case 1:
		//p, _ := strconv.Atoi(codes[0])
		//return &gweb.JsonResult{Data: controller.DistrictService.CityDao.ListC(service.Orm, p)}
	case 2:
		//p, _ := strconv.Atoi(codes[0])
		//cc, _ := strconv.Atoi(codes[1])
		//return &gweb.JsonResult{Data: controller.DistrictService.AreaDao.ListA(service.Orm, p, cc)}

	}
	return &gweb.JsonResult{}
}
