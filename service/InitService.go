package service

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/singleton"
	"log"

	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/tool/encryption"
)

//var GlobalGoodsService = GoodsService{}

func Init(etcd constrain.IEtcd, dbName string) {

	err := singleton.Init(etcd, dbName)
	if err != nil {
		log.Fatalln(err)
	}

	//var err error
	//_db, err := sql.Open("mysql", "tcp:localhost:3306*dandelion/root/123456")
	//_db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/dandelion?charset=utf8mb4&parseTime=True")
	//_db, err := sql.Open("postgres", "postgres://postgres:123456@localhost:5433/dandelion?sslmode=disable")
	//glog.Error(err)
	//dbMap = &gorp.DbMap{Db: _db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8mb4"}}

	//db, err = gorm.Open("mysql", "root:123456@/dandelion?charset=utf8mb4&parseTime=True") //&loc=Local
	//db.SingularTable(true)
	//db = db.Debug()
	//fmt.Println(db)
	//glog.Error(err)
	//defer db.Close()
	//fmt.Println(conf.Config.DBUrl)

	_database := singleton.Orm()

	models := make([]model.IDataBaseFace, 0)

	/*user := &model.User{}
	if _database.HasTable(user) == false {
		_database.Set("gorm:table_options", "AUTO_INCREMENT=1000").CreateTable(user)
	}
	_database.AutoMigrate(user)*/

	models = append(models, &model.User{})
	models = append(models, model.UserInfo{})
	models = append(models, model.Admin{})
	models = append(models, model.Configuration{})
	models = append(models, model.Logger{})
	models = append(models, model.UserJournal{})
	models = append(models, model.StoreJournal{})
	models = append(models, model.CollageRecord{})
	models = append(models, model.CollageGoods{})
	models = append(models, model.District{})
	models = append(models, model.SupplyOrders{})
	models = append(models, model.StoreStock{})
	models = append(models, model.Verification{})
	models = append(models, model.FullCut{})
	models = append(models, model.Voucher{})
	models = append(models, model.TimeSell{})
	models = append(models, model.Store{})
	models = append(models, model.ExpressTemplate{})
	models = append(models, model.CardItem{})
	models = append(models, model.Goods{})
	models = append(models, model.GoodsType{})
	models = append(models, model.GoodsTypeChild{})
	models = append(models, model.OrdersGoods{})
	models = append(models, model.ScoreJournal{})
	models = append(models, model.ScoreGoods{})
	models = append(models, model.Specification{})
	models = append(models, model.Transfers{})
	models = append(models, model.Orders{})
	models = append(models, model.ShoppingCart{})
	models = append(models, model.Rank{})
	models = append(models, model.GiveVoucher{})
	models = append(models, model.Organization{})
	models = append(models, model.OrganizationJournal{})
	models = append(models, model.OrdersPackage{})
	models = append(models, model.Manager{})
	models = append(models, model.Collage{})
	models = append(models, model.Content{})
	models = append(models, model.UserFormIds{})
	models = append(models, model.ContentItem{})
	models = append(models, model.ContentType{})
	models = append(models, model.ContentSubType{})
	models = append(models, model.Question{})
	models = append(models, model.QuestionTag{})
	models = append(models, model.AnswerQuestion{})
	models = append(models, model.TimeSellGoods{})
	models = append(models, model.ContentConfig{})
	models = append(models, model.WXQRCodeParams{})
	models = append(models, model.GoodsAttributes{})
	models = append(models, model.GoodsAttributesGroup{})
	models = append(models, model.GoodsWish{})
	models = append(models, model.LeaveMessage{})
	models = append(models, model.Subscribe{})

	for index := range models {

		if _database.Migrator().HasTable(models[index]) == false {
			//_database.Set("gorm:table_options", "AUTO_INCREMENT=1000").CreateTable(models[index])
			if err := _database.Migrator().CreateTable(models[index]); err != nil {
				glog.Error(err)
			}
		}
		if conf.Config.Debug {
			glog.Debug("migrate:", models[index].TableName())
			if err := _database.AutoMigrate(models[index]); err != nil {
				glog.Error(err)
			}
		}

	}

	var _manager model.Manager
	_database.Where(&model.Manager{Account: "manager"}).First(&_manager)
	if _manager.ID == 0 {
		a := model.Manager{Account: "manager", PassWord: encryption.Md5ByString("274455411")}
		_database.Create(&a)
	}

	//this.Admin.AddAdmin(Name, Password)

	/*var goodsList []model.Goods
	_database.Model(&model.Goods{}).Find(&goodsList)
	for index:=range goodsList{
		goodsItem:=goodsList[index]


		var goodsAtt []model.GoodsAttributes
		json.Unmarshal([]byte(goodsItem.Params),&goodsAtt)

		for iiii:=range goodsAtt{
			attItem:=goodsAtt[iiii]
			glog.Error(goods.GoodsService{}.AddGoodsAttributes(goodsItem.ID,attItem.Name,attItem.Value))
		}


	}*/

	contentTypeList := []model.ContentType{
		{Type: model.ContentTypeContents, Label: "文章列表"},
		{Type: model.ContentTypeContent, Label: "独立文章"},
		{Type: model.ContentTypeIndex, Label: "首页"},
		{Type: model.ContentTypeGallery, Label: "画廊"},
		{Type: model.ContentTypeProducts, Label: "产品"},
		{Type: model.ContentTypeBlog, Label: "博客"},
	}
	for index := range contentTypeList {
		var _contenttype = contentTypeList[index]
		_database.Where(&model.ContentType{Type: _contenttype.Type}).First(&_contenttype)
		if _contenttype.ID == 0 {
			_database.Create(&_contenttype)
		}
	}

}
