package service

import (
	"dandelion/app/service/dao"

	"github.com/nbvghost/gweb/tool"
)

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

	UserInfo := &dao.UserInfo{}
	if _database.HasTable(UserInfo) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=1000").CreateTable(UserInfo)
	}

	admin := &dao.Admin{}
	if _database.HasTable(admin) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(admin)
	}

	item := &dao.Configuration{}
	if _database.HasTable(item) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=1800").CreateTable(item)
	}

	logger := &dao.Logger{}
	if _database.HasTable(logger) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(logger)
	}

	UserJournal := &dao.UserJournal{}
	if _database.HasTable(UserJournal) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(UserJournal)
	}
	StoreJournal := &dao.StoreJournal{}
	if _database.HasTable(StoreJournal) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(StoreJournal)
	}

	District := &dao.District{}
	if _database.HasTable(District) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(District)
	}
	SupplyOrders := &dao.SupplyOrders{}
	if _database.HasTable(SupplyOrders) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(SupplyOrders)
	}
	StoreStock := &dao.StoreStock{}
	if _database.HasTable(StoreStock) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(StoreStock)
	}

	Verification := &dao.Verification{}
	if _database.HasTable(Verification) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Verification)
	}

	FullCut := &dao.FullCut{}
	if _database.HasTable(FullCut) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(FullCut)
	}
	Coupon := &dao.Voucher{}
	if _database.HasTable(Coupon) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Coupon)
	}
	TimeSell := &dao.TimeSell{}
	if _database.HasTable(TimeSell) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(TimeSell)
	}
	Store := &dao.Store{}
	if _database.HasTable(Store) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Store)
	}
	ExpressTemplate := &dao.ExpressTemplate{}
	if _database.HasTable(ExpressTemplate) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(ExpressTemplate)
	}
	CardItem := &dao.CardItem{}
	if _database.HasTable(CardItem) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(CardItem)
	}
	Goods := &dao.Goods{}
	if _database.HasTable(Goods) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Goods)
	}
	GoodsType := &dao.GoodsType{}
	if _database.HasTable(GoodsType) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(GoodsType)
	}
	GoodsTypeChild := &dao.GoodsTypeChild{}
	if _database.HasTable(GoodsTypeChild) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(GoodsTypeChild)
	}
	OrdersGoods := &dao.OrdersGoods{}
	if _database.HasTable(OrdersGoods) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(OrdersGoods)
	}
	ScoreJournal := &dao.ScoreJournal{}
	if _database.HasTable(ScoreJournal) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(ScoreJournal)
	}
	ScoreGoods := &dao.ScoreGoods{}
	if _database.HasTable(ScoreGoods) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(ScoreGoods)
	}
	Specification := &dao.Specification{}
	if _database.HasTable(Specification) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Specification)
	}

	Transfers := &dao.Transfers{}
	if _database.HasTable(Transfers) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Transfers)
	}

	Orders := &dao.Orders{}
	if _database.HasTable(Orders) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Orders)
	}
	ShoppingCart := &dao.ShoppingCart{}
	if _database.HasTable(ShoppingCart) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(ShoppingCart)
	}
	Rank := &dao.Rank{}
	if _database.HasTable(Rank) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Rank)
	}
	GiveCoupon := &dao.GiveVoucher{}
	if _database.HasTable(GiveCoupon) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(GiveCoupon)
	}
	Configuration := &dao.Configuration{}
	if _database.HasTable(Configuration) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Configuration)
	}

	Organization := &dao.Organization{}
	if _database.HasTable(Organization) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Organization)
	}

	/*OrderBrokerageTemp := &dao.OrderBrokerageTemp{}
	if _database.HasTable(OrderBrokerageTemp) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(OrderBrokerageTemp)
	}*/

	OrganizationJournal := &dao.OrganizationJournal{}
	if _database.HasTable(OrganizationJournal) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(OrganizationJournal)
	}

	OrdersPackage := &dao.OrdersPackage{}
	if _database.HasTable(OrdersPackage) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(OrdersPackage)
	}

	Manager := &dao.Manager{}
	if _database.HasTable(Manager) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Manager)
	}

	/*WxConfig := &dao.WxConfig{}
	if _database.HasTable(WxConfig) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(WxConfig)
	}*/

	Article := &dao.Article{}
	if _database.HasTable(Article) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000 CHARSET=utf8mb4").CreateTable(Article)
	}

	Content := &dao.Content{}
	if _database.HasTable(Content) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(Content)
	}
	ContentType := &dao.ContentType{}
	if _database.HasTable(ContentType) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(ContentType)
	}
	ContentSubType := &dao.ContentSubType{}
	if _database.HasTable(ContentSubType) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(ContentSubType)
	}

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

}
