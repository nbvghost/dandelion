package service

import (
	"dandelion/app/play"
	"dandelion/app/service/dao"

	"github.com/nbvghost/gweb/tool"
)

//var GlobalGoodsService = GoodsService{}
var GlobalService GlobalServiceStruct

type GlobalServiceStruct struct {
	Goods GoodsService
}

func init() {

	//var err error
	//_db, err := sql.Open("mysql", "tcp:localhost:3306*dandelion/root/123456")
	//_db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/dandelion?charset=utf8mb4&parseTime=True")
	//_db, err := sql.Open("postgres", "postgres://postgres:123456@localhost:5433/dandelion?sslmode=disable")
	//tool.CheckError(err)
	//dbMap = &gorp.DbMap{Db: _db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8mb4"}}

	//db, err = gorm.Open("mysql", "root:123456@/dandelion?charset=utf8mb4&parseTime=True") //&loc=Local
	//db.SingularTable(true)
	//db = db.Debug()
	//fmt.Println(db)
	//tool.CheckError(err)
	//defer db.Close()
	//fmt.Println(conf.Config.DBUrl)

	_database := dao.Orm()

	user := &dao.User{}
	if _database.HasTable(user) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=1000").CreateTable(user)
	}
	_database.AutoMigrate(user)

	UserInfo := &dao.UserInfo{}
	if _database.HasTable(UserInfo) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=1000").CreateTable(UserInfo)
	}
	_database.AutoMigrate(UserInfo)

	admin := &dao.Admin{}
	if _database.HasTable(admin) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(admin)
	}
	_database.AutoMigrate(admin)

	Configuration := &dao.Configuration{}
	if _database.HasTable(Configuration) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=1800").CreateTable(Configuration)
	}
	_database.AutoMigrate(Configuration)

	logger := &dao.Logger{}
	if _database.HasTable(logger) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(logger)
	}
	_database.AutoMigrate(logger)

	UserJournal := &dao.UserJournal{}
	if _database.HasTable(UserJournal) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(UserJournal)
	}
	_database.AutoMigrate(UserJournal)

	StoreJournal := &dao.StoreJournal{}
	if _database.HasTable(StoreJournal) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(StoreJournal)
	}
	_database.AutoMigrate(StoreJournal)

	District := &dao.District{}
	if _database.HasTable(District) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(District)
	}
	_database.AutoMigrate(District)

	SupplyOrders := &dao.SupplyOrders{}
	if _database.HasTable(SupplyOrders) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(SupplyOrders)
	}
	_database.AutoMigrate(SupplyOrders)

	StoreStock := &dao.StoreStock{}
	if _database.HasTable(StoreStock) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(StoreStock)
	}
	_database.AutoMigrate(StoreStock)

	Verification := &dao.Verification{}
	if _database.HasTable(Verification) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Verification)
	}
	_database.AutoMigrate(Verification)

	FullCut := &dao.FullCut{}
	if _database.HasTable(FullCut) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(FullCut)
	}
	_database.AutoMigrate(FullCut)

	Coupon := &dao.Voucher{}
	if _database.HasTable(Coupon) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Coupon)
	}
	_database.AutoMigrate(Coupon)

	TimeSell := &dao.TimeSell{}
	if _database.HasTable(TimeSell) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(TimeSell)
	}
	_database.AutoMigrate(TimeSell)

	Store := &dao.Store{}
	if _database.HasTable(Store) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Store)
	}
	_database.AutoMigrate(Store)

	ExpressTemplate := &dao.ExpressTemplate{}
	if _database.HasTable(ExpressTemplate) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(ExpressTemplate)
	}
	_database.AutoMigrate(ExpressTemplate)

	CardItem := &dao.CardItem{}
	if _database.HasTable(CardItem) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(CardItem)
	}
	_database.AutoMigrate(CardItem)

	Goods := &dao.Goods{}
	if _database.HasTable(Goods) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Goods)
	}
	_database.AutoMigrate(Goods)

	GoodsType := &dao.GoodsType{}
	if _database.HasTable(GoodsType) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(GoodsType)
	}
	_database.AutoMigrate(GoodsType)

	GoodsTypeChild := &dao.GoodsTypeChild{}
	if _database.HasTable(GoodsTypeChild) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(GoodsTypeChild)
	}
	_database.AutoMigrate(GoodsTypeChild)

	OrdersGoods := &dao.OrdersGoods{}
	if _database.HasTable(OrdersGoods) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(OrdersGoods)
	}
	_database.AutoMigrate(OrdersGoods)

	ScoreJournal := &dao.ScoreJournal{}
	if _database.HasTable(ScoreJournal) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(ScoreJournal)
	}
	_database.AutoMigrate(ScoreJournal)

	ScoreGoods := &dao.ScoreGoods{}
	if _database.HasTable(ScoreGoods) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(ScoreGoods)
	}
	_database.AutoMigrate(ScoreGoods)

	Specification := &dao.Specification{}
	if _database.HasTable(Specification) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Specification)
	}
	_database.AutoMigrate(Specification)

	Transfers := &dao.Transfers{}
	if _database.HasTable(Transfers) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Transfers)
	}
	_database.AutoMigrate(Transfers)

	Orders := &dao.Orders{}
	if _database.HasTable(Orders) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Orders)
	}
	_database.AutoMigrate(Orders)

	ShoppingCart := &dao.ShoppingCart{}
	if _database.HasTable(ShoppingCart) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(ShoppingCart)
	}
	_database.AutoMigrate(ShoppingCart)

	Rank := &dao.Rank{}
	if _database.HasTable(Rank) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Rank)
	}
	_database.AutoMigrate(Rank)

	GiveCoupon := &dao.GiveVoucher{}
	if _database.HasTable(GiveCoupon) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(GiveCoupon)
	}
	_database.AutoMigrate(GiveCoupon)

	Organization := &dao.Organization{}
	if _database.HasTable(Organization) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Organization)
	}
	_database.AutoMigrate(Organization)

	/*OrderBrokerageTemp := &dao.OrderBrokerageTemp{}
	if _database.HasTable(OrderBrokerageTemp) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(OrderBrokerageTemp)
	}*/

	OrganizationJournal := &dao.OrganizationJournal{}
	if _database.HasTable(OrganizationJournal) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(OrganizationJournal)
	}
	_database.AutoMigrate(OrganizationJournal)

	OrdersPackage := &dao.OrdersPackage{}
	if _database.HasTable(OrdersPackage) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(OrdersPackage)
	}
	_database.AutoMigrate(OrdersPackage)

	Manager := &dao.Manager{}
	if _database.HasTable(Manager) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Manager)
	}
	_database.AutoMigrate(Manager)

	Collage := &dao.Collage{}
	if _database.HasTable(Collage) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Collage)
	}
	_database.AutoMigrate(Collage)

	/*WxConfig := &dao.WxConfig{}
	if _database.HasTable(WxConfig) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(WxConfig)
	}*/

	/*tx := _database.Begin()
	for i := 0; i < 700000; i++ {
		user := dao.User{}
		user.Tel = "98321748327492"
		user.Name = tool.CipherEncrypterData("dsfsdfsafsd")
		user.OpenID = tool.CipherEncrypterData("ZxcXZc")
		user.Portrait = tool.CipherEncrypterData("Zxc324324XZc")
		tx.Create(&user)
	}
	tx.Commit()*/

	Article := &dao.Article{}
	if _database.HasTable(Article) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000 CHARSET=utf8mb4").CreateTable(Article)
	}
	_database.AutoMigrate(Article)

	UserFormIds := &dao.UserFormIds{}
	if _database.HasTable(UserFormIds) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000 CHARSET=utf8mb4").CreateTable(UserFormIds)
	}
	_database.AutoMigrate(UserFormIds)

	Content := &dao.Content{}
	if _database.HasTable(Content) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Content)
	}
	_database.AutoMigrate(Content)

	ContentType := &dao.ContentType{}
	if _database.HasTable(ContentType) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(ContentType)
	}
	_database.AutoMigrate(ContentType)

	ContentSubType := &dao.ContentSubType{}
	if _database.HasTable(ContentSubType) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(ContentSubType)
	}
	_database.AutoMigrate(ContentSubType)

	var _manager dao.Manager
	_database.Where(&dao.Manager{Account: "manager"}).First(&_manager)
	if _manager.ID == 0 {
		a := dao.Manager{Account: "manager", PassWord: tool.Md5ByString("274455411")}
		_database.Create(&a)
	}

	//this.Admin.AddAdmin(Name, Password)
	AdminService{}.AddAdmin("admin", "274455411", "")

	var _contenttype dao.ContentType
	_database.Where(&dao.ContentType{Type: "articles"}).First(&_contenttype)
	if _contenttype.ID == 0 {
		a := dao.ContentType{Label: "文章列表", Type: "articles"}
		_database.Create(&a)
	}

	var _Configuration dao.Configuration
	var _configurationService ConfigurationService
	_configurationService.FindWhere(_database, &_Configuration, "K=? and OID=?", play.ConfigurationKey_ScoreConvertGrowValue, 0)
	if _Configuration.ID == 0 {
		a := dao.Configuration{OID: 0, K: play.ConfigurationKey_ScoreConvertGrowValue, V: "1"}
		_configurationService.Add(_database, &a)
	}

}
