package model

import (
	"errors"
	"runtime/debug"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/library/dao"
)

type ConfigurationKey string

const (
	ConfigurationKeyComponentVerifyTicket ConfigurationKey = "ComponentVerifyTicket"
	ConfigurationKeyPoster                ConfigurationKey = "Poster"
	ConfigurationKeyScoreConvertGrowValue ConfigurationKey = "ScoreConvertGrowValue"
	ConfigurationKeyBrokerageLeve1        ConfigurationKey = "BrokerageLeve1"
	ConfigurationKeyBrokerageLeve2        ConfigurationKey = "BrokerageLeve2"
	ConfigurationKeyBrokerageLeve3        ConfigurationKey = "BrokerageLeve3"
	ConfigurationKeyBrokerageLeve4        ConfigurationKey = "BrokerageLeve4"
	ConfigurationKeyBrokerageLeve5        ConfigurationKey = "BrokerageLeve5"
	ConfigurationKeyBrokerageLeve6        ConfigurationKey = "BrokerageLeve6"
	/*
		[{
			      "Matching": [
			        "pages/index/index"
			      ],
			      "Type": "banner",
			      "Images": [
			        {
			          "Src": "/miniapp/images/advert/home/lb1.jpg",
			          "Url": "pages/pro_details/pro_details?ID=2001",
			          "ActonType": "page",
			          "Title": "",
			          "Show": true
			        },
			        {
			          "Src": "/miniapp/images/advert/home/lb2.jpg",
			          "Url": "pages/pro_details/pro_details?ID=2001",
			          "ActonType": "page",
			          "Title": "",
			          "Show": true
			        },
			        {
			          "Src": "/miniapp/images/advert/home/lb3.jpg",
			          "Url": "pages/pro_details/pro_details?ID=2001",
			          "ActonType": "page",
			          "Title": "",
			          "Show": true
			        }
			      ]
			    }]
	*/
	ConfigurationKeyAdvert ConfigurationKey = "Advert"
	/*
		[{
		      "Matching": [
		        "pages/index/index"
		      ],
		      "Type": "pop",
		      "Images": [
		        {
		          "Src": "/miniapp/images/advert/home/pop.png",
		          "Url": "pages/pro_details/pro_details?ID=2001",
		          "ActonType": "page",
		          "Title": "",
		          "Show": false
		        }
		      ]
		    }]
	*/
	ConfigurationKeyPop ConfigurationKey = "Pop"
	/*
		[
		    {
		      "Src": "/miniapp/images/icon/index/mzf.png",
		      "Url": "",
		      "ActonType": "webview",
		      "Title": "微信支付",
		      "Show": true
		    },
		    {
		      "Src": "/miniapp/images/icon/index/mlg.png",
		      "Url": "",
		      "ActonType": "webview",
		      "Title": "担保交易",
		      "Show": true
		    },
		    {
		      "Src": "/miniapp/images/icon/index/mlh.png",
		      "Url": "",
		      "ActonType": "webview",
		      "Title": "优先商品",
		      "Show": true
		    },
		    {
		      "Src": "/miniapp/images/icon/index/mls.png",
		      "Url": "",
		      "ActonType": "webview",
		      "Title": "全场正品",
		      "Show": true
		    },
		    {
		      "Src": "/miniapp/images/icon/index/msh.png",
		      "Url": "",
		      "ActonType": "webview",
		      "Title": "快捷售后",
		      "Show": true
		    },
		    {
		      "Src": "/miniapp/images/icon/index/msh.png",
		      "Url": "https://web.sites.ink/make-shape/index",
		      "ActonType": "webview",
		      "Title": "2D",
		      "Show": true
		    },
		    {
		      "Src": "/miniapp/images/icon/index/msh.png",
		      "Url": "https://web.sites.ink/make-shape/3d",
		      "ActonType": "webview",
		      "Title": "3D",
		      "Show": true
		    }
		  ]
	*/
	ConfigurationKeyQuickLink              ConfigurationKey = "QuickLink"
	ConfigurationKeyPaymentPaypalClientId  ConfigurationKey = "PaymentPaypalClientId"
	ConfigurationKeyPaymentPaypalAppSecret ConfigurationKey = "PaymentPaypalAppSecret"
)

type Configuration struct {
	dao.Entity
	OID dao.PrimaryKey   `gorm:"column:OID"`
	K   ConfigurationKey `gorm:"column:K"`
	V   string           `gorm:"column:V"`
}

func (u *Configuration) BeforeCreate(scope *gorm.DB) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))

	}
	return nil
}
func (Configuration) TableName() string {
	return "Configuration"
}
