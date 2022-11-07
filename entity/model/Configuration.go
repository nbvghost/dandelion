package model

import (
	"errors"
	"runtime/debug"

	"gorm.io/gorm"

	"github.com/nbvghost/gpa/types"
)

type ConfigurationKey uint

const (
	ConfigurationKeyComponentVerifyTicket ConfigurationKey = 1001
	ConfigurationKeyPoster                ConfigurationKey = 1002
	ConfigurationKeyScoreConvertGrowValue ConfigurationKey = 1100
	ConfigurationKeyBrokerageLeve1        ConfigurationKey = 1201
	ConfigurationKeyBrokerageLeve2        ConfigurationKey = 1202
	ConfigurationKeyBrokerageLeve3        ConfigurationKey = 1203
	ConfigurationKeyBrokerageLeve4        ConfigurationKey = 1204
	ConfigurationKeyBrokerageLeve5        ConfigurationKey = 1205
	ConfigurationKeyBrokerageLeve6        ConfigurationKey = 1206
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
	ConfigurationKeyAdvert ConfigurationKey = 1300
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
	ConfigurationKeyPop ConfigurationKey = 1301
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
	ConfigurationKeyQuickLink ConfigurationKey = 1302
)

type Configuration struct {
	types.Entity
	OID types.PrimaryKey `gorm:"column:OID"`
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
