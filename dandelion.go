package main

import (
	"github.com/nbvghost/dandelion/app/action/account"
	"github.com/nbvghost/dandelion/app/action/admin"
	"github.com/nbvghost/dandelion/app/action/api"
	"github.com/nbvghost/dandelion/app/action/file"
	"github.com/nbvghost/dandelion/app/action/images"
	"github.com/nbvghost/dandelion/app/action/index"
	"github.com/nbvghost/dandelion/app/action/manager"
	"github.com/nbvghost/dandelion/app/action/mp"
	"github.com/nbvghost/dandelion/app/action/payment"
	"github.com/nbvghost/dandelion/app/action/sites"
	"github.com/nbvghost/dandelion/app/service"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/conf"
	"net/http"
)

func init() {
	glog.Param.PushAddr = conf.Config.LogServer
	glog.Param.AppName = "dandelion"
	glog.Param.LogFilePath = conf.Config.LogDir
	glog.Param.StandardOut = true
	glog.Param.FileStorage = false
	glog.Start()
}

func main() {

	service.Init()

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
	sites.NewController("/sites/", sites)

	file := &file.Controller{}
	file.NewController("/file", file)

	_http := &http.Server{
		Addr:    conf.Config.HttpPort,
		Handler: nil,
	}
	_https := &http.Server{
		Addr:    conf.Config.HttpsPort,
		Handler: nil,
	}
	gweb.StartServer(http.DefaultServeMux, _http, _https)
}
