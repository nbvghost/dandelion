package service

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service/internal/activity"
	"github.com/nbvghost/dandelion/service/internal/admin"
	"github.com/nbvghost/dandelion/service/internal/cache"
	"github.com/nbvghost/dandelion/service/internal/catch"
	"github.com/nbvghost/dandelion/service/internal/company"
	"github.com/nbvghost/dandelion/service/internal/configuration"
	"github.com/nbvghost/dandelion/service/internal/content"
	"github.com/nbvghost/dandelion/service/internal/express"
	"github.com/nbvghost/dandelion/service/internal/file"
	"github.com/nbvghost/dandelion/service/internal/goods"
	"github.com/nbvghost/dandelion/service/internal/journal"
	"github.com/nbvghost/dandelion/service/internal/logger"
	"github.com/nbvghost/dandelion/service/internal/manager"
	"github.com/nbvghost/dandelion/service/internal/network"
	"github.com/nbvghost/dandelion/service/internal/order"
	"github.com/nbvghost/dandelion/service/internal/pinyin"
	"github.com/nbvghost/dandelion/service/internal/question"
	"github.com/nbvghost/dandelion/service/internal/search"
	"github.com/nbvghost/dandelion/service/internal/site"
	"github.com/nbvghost/dandelion/service/internal/sms"
	"github.com/nbvghost/dandelion/service/internal/task"
	"github.com/nbvghost/dandelion/service/internal/user"
	"github.com/nbvghost/dandelion/service/internal/wechat"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

var Content = content.ContentService{}
var Activity = struct {
	CardItem    activity.CardItemService
	Collage     activity.CollageService
	FullCut     activity.FullCutService
	GiveVoucher activity.GiveVoucherService
	Rank        activity.RankService
	ScoreGoods  activity.ScoreGoodsService
	Settlement  activity.SettlementService
	TimeSell    activity.TimeSellService
	Voucher     activity.VoucherService
}{}
var Admin = admin.AdminService{}
var Catch = catch.SpiderService{}
var Company = struct {
	Organization company.OrganizationService
	Store        company.StoreService
}{}
var Configuration = configuration.ConfigurationService{}
var Express = struct {
	ExpressTemplate express.ExpressTemplateService
	District        express.DistrictService
}{}
var File = struct {
	File file.FileService
	Html file.HtmlService
}{}
var Goods = struct {
	Goods         goods.GoodsService
	Attributes    goods.AttributesService
	GoodsType     goods.GoodsTypeService
	SKU           goods.SKUService
	Sort          goods.SortService
	Specification goods.SpecificationService
	Tag           goods.TagService
	//ProductOptions func(ctx constrain.IContext, oid dao.PrimaryKey) (*serviceargument.Options, error)
}{
	//ProductOptions: goods.ProductOptions,
}

var Journal = struct {
	journal.JournalService
	NewDataTypeUser func(UserID dao.PrimaryKey) journal.IDataType
}{
	NewDataTypeUser: journal.NewDataTypeUser,
}
var Logger = logger.LoggerService{}
var Manager = manager.ManagerService{}
var Order = struct {
	Orders       order.OrdersService
	ShoppingCart order.ShoppingCartService
	Transfers    order.TransfersService
	Verification order.VerificationService
}{}
var Pinyin = pinyin.Service{}
var Question = question.QuestionService{}
var Search = search.Service{}
var Site = site.Service{}
var SMS = sms.Service{}
var Task = task.TimeTaskService{}
var User = user.UserService{}
var Wechat = struct {
	Wx             wechat.WxService
	WXQRCodeParams wechat.WXQRCodeParamsService
	MessageNotify  wechat.MessageNotify
}{}
var Network = struct {
	SMS    network.SMS
	Email  network.Email
	NewSMS func(oid dao.PrimaryKey) *network.SMS
}{
	NewSMS: network.NewSMS,
}

var Cache = struct {
	ChinesePinyinCache cache.ChinesePinyinCache
	LanguageCache      cache.LanguageCache
	LanguageCodeCache  cache.LanguageCodeCache
	RedisCache         cache.RedisCache
}{
	ChinesePinyinCache: cache.ChinesePinyinCache{Pinyin: make(map[string]string)},
	LanguageCache:      cache.LanguageCache{ShowLanguage: make([]model.Language, 0)},
	LanguageCodeCache:  cache.LanguageCodeCache{LangBaiduCode: make(map[string]string)},
}

func init() {

}

func GetSiteData[T serviceargument.ListType](context constrain.IContext, OID dao.PrimaryKey) serviceargument.SiteData[T] {

	var moduleContentData serviceargument.SiteData[T]

	var item model.ContentItem
	var subItem = model.ContentSubType{Uri: "all"}

	currentMenuData := serviceargument.NewMenusData(item, subItem)

	menusData := Site.FindShowMenus(OID)
	for _, v := range menusData.List {
		if v.ID == currentMenuData.TypeID {
			currentMenuData.Menus = v
			break
		}
	}

	contentItemMap := repository.ContentItemDao.ListContentItemByOIDMap(OID)

	allMenusData := Site.FindAllMenus(OID)

	tags := Content.FindContentTagsByContentItemID(OID, currentMenuData.TypeID)

	var navigations []extends.Menus

	var typeNameMap = make(map[dao.PrimaryKey]extends.Menus)

	for index, v := range menusData.List {
		if v.ID == currentMenuData.TypeID {
			navigations = append(navigations, menusData.List[index])
			for si, sv := range v.List {
				typeNameMap[sv.ID] = sv
				if sv.ID == currentMenuData.SubTypeID {
					navigations = append(navigations, menusData.List[index].List[si])
				} else {
					for _, ssv := range sv.List {
						typeNameMap[ssv.ID] = ssv
					}
					for ssi, ssv := range sv.List {
						if ssv.ID == currentMenuData.SubTypeID {
							navigations = append(navigations, menusData.List[index].List[si])
							navigations = append(navigations, menusData.List[index].List[si].List[ssi])
							break
						}
					}
				}
			}
			break
		}
	}

	organization := Company.Organization.GetOrganization(OID).(*model.Organization)
	contentConfig := repository.ContentConfigDao.GetContentConfig(db.Orm(), OID)

	menusPage := allMenusData.ListMenusByType(model.ContentTypePage)
	moduleContentData = serviceargument.SiteData[T]{
		AllMenusData:    allMenusData,
		MenusData:       menusData,
		PageMenus:       menusPage,
		CurrentMenuData: currentMenuData,
		ContentItem:     item,
		ContentSubType:  subItem,
		Pagination:      serviceargument.Pagination[T]{},
		Tags:            tags,
		Navigations:     navigations,
		Organization:    *organization,
		ContentConfig:   contentConfig,
		TypeNameMap:     typeNameMap,
		ContentItemMap:  contentItemMap,
	}

	companyName := contentConfig.Name
	if len(companyName) == 0 {
		companyName = organization.Name
	}
	moduleContentData.SiteAuthor = companyName

	return moduleContentData
}
