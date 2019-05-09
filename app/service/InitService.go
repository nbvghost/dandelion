package service

import (
	"dandelion/app/service/dao"

	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/gweb/tool"
)

//var GlobalGoodsService = GoodsService{}
var GlobalService GlobalServiceStruct

type GlobalServiceStruct struct {
	Goods  GoodsService
	Orders OrdersService
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

	models := make([]dao.IDataBaseFace, 0)

	/*user := &dao.User{}
	if _database.HasTable(user) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=1000").CreateTable(user)
	}
	_database.AutoMigrate(user)*/

	models = append(models, dao.User{})
	models = append(models, dao.UserInfo{})
	models = append(models, dao.Admin{})
	models = append(models, dao.Configuration{})
	models = append(models, dao.Logger{})
	models = append(models, dao.UserJournal{})
	models = append(models, dao.StoreJournal{})
	models = append(models, dao.CollageRecord{})
	models = append(models, dao.District{})
	models = append(models, dao.SupplyOrders{})
	models = append(models, dao.StoreStock{})
	models = append(models, dao.Verification{})
	models = append(models, dao.FullCut{})
	models = append(models, dao.Voucher{})
	models = append(models, dao.TimeSell{})
	models = append(models, dao.Store{})
	models = append(models, dao.ExpressTemplate{})
	models = append(models, dao.CardItem{})
	models = append(models, dao.Goods{})
	models = append(models, dao.GoodsType{})
	models = append(models, dao.GoodsTypeChild{})
	models = append(models, dao.OrdersGoods{})
	models = append(models, dao.ScoreJournal{})
	models = append(models, dao.ScoreGoods{})
	models = append(models, dao.Specification{})
	models = append(models, dao.Transfers{})
	models = append(models, dao.Orders{})
	models = append(models, dao.ShoppingCart{})
	models = append(models, dao.Rank{})
	models = append(models, dao.GiveVoucher{})
	models = append(models, dao.Organization{})
	models = append(models, dao.OrganizationJournal{})
	models = append(models, dao.OrdersPackage{})
	models = append(models, dao.Manager{})
	models = append(models, dao.Collage{})
	models = append(models, dao.Article{})
	models = append(models, dao.UserFormIds{})
	models = append(models, dao.Content{})
	models = append(models, dao.ContentType{})
	models = append(models, dao.ContentSubType{})
	models = append(models, dao.Question{})
	models = append(models, dao.QuestionTag{})
	models = append(models, dao.AnswerQuestion{})

	for index := range models {

		if _database.HasTable(models[index]) == false {
			//_database.Set("gorm:table_options", "AUTO_INCREMENT=1000").CreateTable(models[index])
			_database.CreateTable(models[index])
		}
		if conf.Config.Debug {
			_database.AutoMigrate(models[index])
		}

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
