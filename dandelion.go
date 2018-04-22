package main

import (
	"dandelion/app/action/account"
	"dandelion/app/action/admin"
	"dandelion/app/action/district"
	"dandelion/app/action/file"
	"dandelion/app/action/front"
	"dandelion/app/action/project"

	"dandelion/app/action/index"
	"dandelion/app/action/manager"
	"dandelion/app/action/order"
	"dandelion/app/action/wx"
	"dandelion/app/util"
	"fmt"
	"log"
	"net/http"

	"github.com/nbvghost/gweb/conf"
)

func init() {
	index := &index.Controller{}
	index.NewController("/", index)
	//fmt.Println(index)

	front := &front.Controller{}
	front.NewController("/front/", front)

	account := &account.Controller{}
	account.NewController("/account/", account)

	manager := &manager.Controller{}
	manager.NewController("/manager/", manager)

	file := &file.Controller{}
	file.NewController("/file/", file)

	order := &order.Controller{}
	order.NewController("/order/", order)

	wx := &wx.Controller{}
	wx.NewController("/wx/", wx)

	district := &district.Controller{}
	district.NewController("/district/", district)

	admin := &admin.Controller{}
	admin.NewController("/admin/", admin)

	project := &project.Controller{}
	project.NewController("/project/", project)

	fmt.Println(util.StructToMap(admin))

	//http.Handle((&admin.Controller{}).Init())
	//http.Handle((&account.Controller{}).Init())
	//http.Handle((&manager.Controller{}).Init())
	//http.Handle((&file.Controller{}).Init())

	//http.Handle((&order.Controller{}).Init())
	//http.Handle((&wx.Controller{}).Init())*/

}
func main() {

	//encrypt_type=aes&msg_signature=50e021cd61ccc9b0f8e04c73027b91a379286ad2&nonce=869164555&signature=cdfbb8836edeed630c28d66163e56d7d16248f96&timestamp=1510646867

	//fmt.Println(fdsf)

	/*toke:=&wxpay.TokenXML{}
	xml.Unmarshal([]byte(fdsf),toke)
	fmt.Println(toke)


	sdfd,content:=wxpay.DecryptMsg("50e021cd61ccc9b0f8e04c73027b91a379286ad2","1510646867","869164555",toke.Encrypt)
	fmt.Println(sdfd)
	fmt.Println(content)



	sdfd,content=wxpay.DecryptMsg("50e021cd61ccc9b0f8e04c73027b91a379286ad2","1510646867","869164555","ticket@@@W1Dee2YsQJyfYYyeuo7T9woqDFpXexdnfzM7_X_Egfz-TqDOKU4OAgNFdLbQnBrEsR452l20jK7bMgkZQYdLFg")
	fmt.Println(sdfd)
	fmt.Println(content)*/

	//fmt.Println(wxpay.Decrypt("ticket@@@W1Dee2YsQJyfYYyeuo7T9woqDFpXexdnfzM7_X_Egfz-TqDOKU4OAgNFdLbQnBrEsR452l20jK7bMgkZQYdLFg"))

	/*var kks dao.User

	asd := reflect.TypeOf(kks)
	fmt.Println(asd.Kind())
	fmt.Println(asd.FieldByName("ID"))
	fmt.Println(asd)

	reg := regexp.MustCompile("(/){2,}")
	fmt.Println(reg.ReplaceAllString("//", "/"))

	fmt.Println(strings.Replace("/////////////", "//", "/", -1))*/

	//go service.ReadANetArticle()
	//go service.ReadBNetArticle()
	//service.ReadWeiXinArticle("http://mp.weixin.qq.com/s?src=3&timestamp=1502797326&ver=1&signature=ZkWvEzc20NzLSR5Z-kzfaLLpHYdlUqEmkhRt*Lt-2ZMxv9-*ymqMNwg1INMM56CAI1Psx0QeXe6L2woyoS7f3W9BDd5Vm1GtSrDtl1FnfdfegToMLUqITz8AbDj1BH09fyj1xAz5htcRvUCqeHT8CktsIMXRg1q6T7vvrVzvf28=")

	//Path := "/////sdfdsf/:a/:b/:c/:d/sfds/////"
	//CuPath := "////sdfdsf/a/b/c/d/sfds/////" ///sdfdsf/d:idf/d{sfsdf}f/{fdgfd}/fd/sfds
	//fmt.Println(tools.MatchURL("/////sdfdsf/:a/:b/:c/:d/sfds/////", "////sdfdsf/a/b/c/d/sfds/////"))
	//tools.Trace("web server listen to 9000")

	/*_, err := os.Open("cert/server.crt")
	if err != nil {
		panic(err)
	}*/

	/*f, _ := os.Open("resources/geo/js.txt")

	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		if !strings.EqualFold(line, "") {
			stk := strings.Split(line, " ")

			code := stk[0]
			name := strings.TrimSpace(stk[len(stk)-1])

			ID, _ := strconv.ParseUint(code, 10, 64)

			p, _ := strconv.Atoi(code[:2])
			c, _ := strconv.Atoi(code[2:4])
			a, _ := strconv.Atoi(code[4:6])

			fmt.Println(code, name)
			fmt.Println(p, c, a)

			if c == 0 && a == 0 {
				service.District.Add(&dao.Province{ID: ID, P: p, Name: name})
			} else if a == 0 {
				service.District.Add(&dao.City{ID: ID, P: p, C: c, Name: name})
			} else {
				if strings.EqualFold(name, "市辖区") {
					continue
				}
				service.District.Add(&dao.Area{ID: ID, P: p, C: c, A: a, Name: name})
			}
		}

		if err != nil {
			if err == io.EOF {
				fmt.Println(err)
				break
			}
			fmt.Println(err)
		}
	}*/

	go func() {

		/*pool := x509.NewCertPool()
		caCertPath := "cert/1_root_bundle.crt"

		caCrt, err := ioutil.ReadFile(caCertPath)
		if err != nil {
			fmt.Println("ReadFile err:", err)
			return
		}
		pool.AppendCertsFromPEM(caCrt)

		s := &http.Server{
			Addr:      ":443",
			Handler:   nil,
			TLSConfig: &tls.Config{
			//Certificates:[]tls.Certificate{}{cc}
			//RootCAs: pool,
			//ClientAuth: tls.RequireAndVerifyClientCert,
			//ClientAuth: tls.NoClientCert,
			},
		}*/

		http.ListenAndServeTLS(conf.Config.HttpsPort, conf.Config.TLSCertFile, conf.Config.TLSKeyFile, nil)

	}()

	err := http.ListenAndServe(conf.Config.HttpPort, nil)
	log.Println(err)
	//http.ListenAndServeTLS(":9000", "cert/server.crt", "cert/server.key", nil)
}
