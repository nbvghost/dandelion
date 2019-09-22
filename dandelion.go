package main

import (
	"net/http"

	"github.com/nbvghost/dandelion/app/action/sites"
	"github.com/nbvghost/dandelion/app/service"

	"github.com/nbvghost/gweb"

	"github.com/nbvghost/dandelion/app/action/account"
	"github.com/nbvghost/dandelion/app/action/admin"
	"github.com/nbvghost/dandelion/app/action/api"
	"github.com/nbvghost/dandelion/app/action/images"
	"github.com/nbvghost/dandelion/app/action/index"
	"github.com/nbvghost/dandelion/app/action/manager"
	"github.com/nbvghost/dandelion/app/action/mp"
	"github.com/nbvghost/dandelion/app/action/payment"

	"github.com/nbvghost/gweb/conf"
)

func main() {

	s := &service.Catch1688Service{}
	s.Run()

	//service.WxService{}.SendUniformMessage(service.WxService{}.MiniWeb(), service.WxService{}.MiniProgram())

	//sms:=service.SMSService{}
	//sms.SendAliyunSms(map[string]interface{}{"name":1,"usename":454,"sum":45,"time":45,"ordernum":544,"productname":5454},"SMS_137687089","15959898368","LTAIqUS1pQEIx6Gr","Oekyziw358oTQPVlbrL1IMzlDFV4ce")
	//sms.SendAliyunSms(map[string]interface{}{"customer":"test"},"SMS_71390007","15300000001","LTAIqUS1pQEIx6Gr","Oekyziw358oTQPVlbrL1IMzlDFV4ce")

	// go func() {
	// 	for {
	// 		resp, err := http.Get("http://api.m.taobao.com/rest/api3.do?api=mtop.common.getTimestamp")
	// 		if err == nil {
	// 			//fmt.Println(err)
	// 			//fmt.Println(resp)
	// 			b, _ := ioutil.ReadAll(resp.Body)
	// 			result := make(map[string]interface{})
	// 			json.Unmarshal(b, &result)
	// 			if result["data"] != nil {

	// 				if result["data"].(map[string]interface{})["t"] != nil {

	// 					_now, _ := strconv.ParseUint(result["data"].(map[string]interface{})["t"].(string), 10, 64)
	// 					now := time.Unix(int64(_now/1000), 0)
	// 					//outtime,_:=time.Parse("2006-01-02 15:04:05","2018-06-01 00:00:00")
	// 					outtime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2018-05-20 00:00:00", now.Location())
	// 					if now.Unix() > outtime.Unix() {
	// 						os.Exit(0)
	// 					}
	// 				}
	// 			}
	// 		}

	// 		time.Sleep(time.Minute * 3)
	// 	}
	// }()

	/*users:=service.UserService{}
	u:=users.FindUserByOpenID(dao.Orm(),"oy4XD5BNScPCqH2Jkan4NBfYhjDA")
	for ia:=0;ia<15;ia++{

		mu:=dao.User{Name:"一级",Tel:"2445326445"+strconv.Itoa(ia),SuperiorID:u.ID}
		users.Add(dao.Orm(),&mu)

		for ib:=0;ib<15;ib++{

			mu:=dao.User{Name:"二级",Tel:"2445326445"+strconv.Itoa(ia)+strconv.Itoa(ib),SuperiorID:mu.ID}
			users.Add(dao.Orm(),&mu)

			for ic:=0;ic<15;ic++{

				mu:=dao.User{Name:"三级",Tel:"2445326445"+strconv.Itoa(ia)+strconv.Itoa(ib)+strconv.Itoa(ic),SuperiorID:mu.ID}
				users.Add(dao.Orm(),&mu)

				for id:=0;id<15;id++{

					mu:=dao.User{Name:"四级",Tel:"2445326445"+strconv.Itoa(ia)+strconv.Itoa(ib)+strconv.Itoa(ic)+strconv.Itoa(id),SuperiorID:mu.ID}
					users.Add(dao.Orm(),&mu)


					for ie:=0;ie<15;ie++{

						mu:=dao.User{Name:"五级",Tel:"2445326445"+strconv.Itoa(ia)+strconv.Itoa(ib)+strconv.Itoa(ic)+strconv.Itoa(id)+strconv.Itoa(ie),SuperiorID:mu.ID}
						users.Add(dao.Orm(),&mu)

						for ig:=0;ig<15;ig++{

							mu:=dao.User{Name:"六级",Tel:"2445326445"+strconv.Itoa(ia)+strconv.Itoa(ib)+strconv.Itoa(ic)+strconv.Itoa(id)+strconv.Itoa(ie)+strconv.Itoa(ig),SuperiorID:mu.ID}
							users.Add(dao.Orm(),&mu)

						}

					}
				}

			}

		}

	}
	fmt.Println(u)*/

	/*sersdf := service.JournalService{}
	fmt.Println(sersdf.StoreListJournal(2000, "2018-4-27", "2018-4-27"))
	if true {
		return
	}*/

	//su, d := wxpay.Decrypt("3JjPMHj4548SCA9dmo+ogUrGSSu3x2fUOH5AZnacuis37txpeyVfn5sbVfgIiOGQ5+XU/FHhSM9cIuw0zS76/p69lsqQoxW7TiQTQS/yzjuLyYbRTRvuXbnxBk4fTysCwx6ITdVFIeUhFtgitKTllFbk1JMGzUWmLKbeHBRPWRZIxL754ykN6X9x3A7ydXZ4tqAzyqjq+ZEzHn9tl7apLjbU1oyH4c3mt1YQSuVgqK1enyksov7q6+gTAsKsvW0ftgjReLRrk8+KiS7OFiOgXlyHoHrro2LR2BgB6LaZqtdGKYJwA6wWCCBAiFYZjqWG7/Bs3wW0L9wnVV8pD5Kc9768iXetspKfcPGFoCI7xu3xvI/0KFxBBhk75ADEPtwVMGnPtp62TiCFWxyEBX+nPkCT0OlRDfyN/JTb0VVnMs1pqwQe2NogjnGNEQLI2sNFNgrGOuWIbw+aZILStoQV+Q==", "E9jPA7zOICjC0/Ldt/zBmA==", "jnA5dg8DE32U8gts7lXAAg==")
	//fmt.Println(su, d)

	//front := &front.Controller{}
	//front.NewController("/front/", front)

	admin := &admin.Controller{}
	admin.NewController("/admin/", admin)

	manager := &manager.Controller{}
	manager.NewController("/manager/", manager)

	account := &account.Controller{}
	account.NewController("/account/", account)

	images := &images.Controller{}
	images.NewController("/images/", images)

	mp := &mp.Controller{}
	mp.NewController("/mp/", mp)

	payment := &payment.Controller{}
	payment.NewController("/payment/", payment)

	home := &index.Controller{}
	home.NewController("/", home)

	api := &api.Controller{}
	api.NewController("/api", api)

	sites := &sites.Controller{}
	sites.NewController("/sites", sites)

	_http := &http.Server{
		Addr:    conf.Config.HttpPort,
		Handler: nil,
	}
	_https := &http.Server{
		Addr:    conf.Config.HttpsPort,
		Handler: nil,
	}
	_https = nil
	gweb.StartServer(http.DefaultServeMux, _http, _https)
}
