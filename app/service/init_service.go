package service

import (
	"time"

	"dandelion/app/service/dao"

	"github.com/jinzhu/gorm"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/gweb/tool"
)

var (
	Admin       = AdminService{}
	Appointment = AppointmentService{}
	Article     = ArticleService{}
	Category    = CategoryService{}
	File        = FileService{}
	Html        = HtmlService{}
	Manager     = ManagerService{}
	Company     = CompanyService{}
	User        = UserService{}
	Classify    = ClassifyService{}
	OrderPack   = OrderPackService{}

	Configuration = ConfigurationService{}
	WxConfig      = WxConfigService{}
	District      = DistrictService{}
	Orm           *gorm.DB
)

func init() {
	var err error
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

	//root:123456@tcp(127.0.0.1:3306)/dandelion?charset=utf8&parseTime=True&loc=Local
	Orm, err = gorm.Open("mysql", conf.Config.DBUrl)
	tool.CheckError(err)

	Orm.Debug()
	Orm.LogMode(true)

	Orm.Exec("SET GLOBAL GROUP_CONCAT_MAX_LEN=1844674407370954752")
	Orm.Exec("SET SESSION GROUP_CONCAT_MAX_LEN=1844674407370954752")

	user := &dao.User{}
	if Orm.HasTable(user) == false {
		Orm.Set("gorm:table_options", "AUTO_INCREMENT=1000").CreateTable(user)
	}

	Project := &dao.Project{}
	if Orm.HasTable(Project) == false {
		Orm.Set("gorm:table_options", "AUTO_INCREMENT=1000").CreateTable(Project)
	}

	admin := &dao.Admin{}
	if Orm.HasTable(admin) == false {
		Orm.Set("gorm:table_options", "AUTO_INCREMENT=2000").CreateTable(admin)
	}

	company := &dao.Company{}
	if Orm.HasTable(company) == false {
		Orm.Set("gorm:table_options", "AUTO_INCREMENT=3000").CreateTable(company)
	}

	manager := &dao.Manager{}
	if Orm.HasTable(manager) == false {
		Orm.Set("gorm:table_options", "AUTO_INCREMENT=4000").CreateTable(manager)
	}

	article := &dao.Article{}
	if Orm.HasTable(article) == false {
		Orm.Set("gorm:table_options", "AUTO_INCREMENT=5000").CreateTable(article)
	}

	category := &dao.Category{}
	if Orm.HasTable(category) == false {
		Orm.Set("gorm:table_options", "AUTO_INCREMENT=6000").CreateTable(category)
	}
	classify := &dao.Classify{}
	if Orm.HasTable(classify) == false {
		Orm.Set("gorm:table_options", "AUTO_INCREMENT=7000").CreateTable(classify)
	}
	appointment := &dao.Appointment{}
	if Orm.HasTable(appointment) == false {
		Orm.Set("gorm:table_options", "AUTO_INCREMENT=8000").CreateTable(appointment)
	}

	orderPack := &dao.OrderPack{}
	if Orm.HasTable(orderPack) == false {
		Orm.Set("gorm:table_options", "AUTO_INCREMENT=9000").CreateTable(orderPack)
	}

	shopCoupon := &dao.ShopCoupon{}
	if Orm.HasTable(shopCoupon) == false {
		Orm.Set("gorm:table_options", "AUTO_INCREMENT=1000").CreateTable(shopCoupon)
	}

	userCoupon := &dao.UserCoupon{}
	if Orm.HasTable(userCoupon) == false {
		Orm.Set("gorm:table_options", "AUTO_INCREMENT=1100").CreateTable(userCoupon)
	}

	wxConfig := &dao.WxConfig{}
	if Orm.HasTable(wxConfig) == false {
		Orm.Set("gorm:table_options", "AUTO_INCREMENT=1700").CreateTable(wxConfig)
	}

	item := &dao.Configuration{}
	if Orm.HasTable(item) == false {
		Orm.Set("gorm:table_options", "AUTO_INCREMENT=1800").CreateTable(item)
	}

	province := &dao.Province{}
	if Orm.HasTable(province) == false {
		Orm.CreateTable(province)
	}

	city := &dao.City{}
	if Orm.HasTable(city) == false {
		Orm.CreateTable(city)
	}

	area := &dao.Area{}
	if Orm.HasTable(area) == false {
		Orm.CreateTable(area)
	}

	var _manager dao.Manager

	Orm.Where(&dao.Manager{Account: "manager"}).First(&_manager)

	if _manager.ID == 0 {
		/*admin := &Admin{}
		admin.Accounts = "admin"
		admin.LoginPassWord = "E10ADC3949BA59ABBE56E057F20F883E"*/
		//WriteAdmin(admin)

		a := dao.Manager{Account: "manager", PassWord: "E10ADC3949BA59ABBE56E057F20F883E"}
		a.CreatedAt = time.Now()
		//orm.NewRecord(a) // => returns `true` as primary key is blank
		Orm.Create(&a)
		//orm.NewRecord(a) // => return `false` after `user` created
		//adminCollection.Find()
	}

	/*article := dbMap.AddTable(Article{})
	article.SetKeys(true, "ID")
	article.ColMap("Content").SetMaxSize(math.MaxInt64)

	err = dbMap.CreateTablesIfNotExists()
	tool.CheckError(err)

	dbMap.Insert(&User{Name: `lsdkfjsa`, Password: tools.Md5(`274455411`), Email: `nbvghost@qq.com`, Tel: `13809549424`, Authority: ManagerType, LastLoginAt: time.Time{}, CreatedAt: time.Now()})

	ff := &User{}
	err = dbMap.SelectOne(ff, "select * from User where Email=?", "nbvghost@qq.com")
	tool.CheckError(err)
	fmt.Println(ff)*/
	//db.Create(&User{Name: `lsdkfjsa`, Password: tools.Md5(`274455411`), Email: `nbvghost@qq.com`, Tel: `13809549424`, Authority: ManagerType, Last_Login_At: time.Time{}, Created_At: time.Now()})
	Admin.AdminDao.FindAdmin(Orm)
}
