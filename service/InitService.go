package service

import (
	"fmt"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/constrain/key"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/environments"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/service/cache"

	"github.com/nbvghost/glog"
	"github.com/nbvghost/tool/encryption"
)

func Init(app key.MicroServer, etcd constrain.IEtcd, dbName string) error {
	err := singleton.Init(etcd, dbName)
	if err != nil {
		return err
	}

	_database := singleton.Orm()

	models := make([]model.IDataBaseFace, 0)

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
	models = append(models, model.FullTextSearch{})
	models = append(models, model.Pinyin{})
	models = append(models, model.Language{})
	models = append(models, model.Translate{})
	models = append(models, model.DNS{})
	models = append(models, model.Advert{})
	models = append(models, model.WechatConfig{})
	models = append(models, model.PushData{})

	//set db session application name
	_database.Exec(fmt.Sprintf("SET application_name='%s'", app))

	dbContentFunc := `CREATE OR REPLACE FUNCTION process_content_full_text_search() RETURNS TRIGGER AS
$Content$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        Delete from "FullTextSearch" where "TID" = OLD."ID" and "Type" = 'content';
    ELSIF (TG_OP = 'UPDATE') THEN
        update "FullTextSearch"
        set "UpdatedAt"=NEW."UpdatedAt",
            "Title"=NEW."Title",
            "Content"=NEW."Content",
            "Picture"=NEW."Picture",
            "Index"=setweight(to_tsvector('english', coalesce(NEW."Title", '')),'A') || setweight(to_tsvector('english', coalesce(NEW."Content", '')),'B')
        where "TID" = NEW."ID"
          and "Type" = 'content';
    ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO "FullTextSearch" ("ID", "CreatedAt", "UpdatedAt", "OID", "TID", "Title", "Content", "Picture",
                                      "Type", "Index")
        values (DEFAULT, NEW."CreatedAt", NEW."UpdatedAt", NEW."OID", NEW."ID", NEW."Title", NEW."Content",
                NEW."Picture", 'content',
                setweight(to_tsvector('english', coalesce(NEW."Title", '')),'A') || setweight(to_tsvector('english', coalesce(NEW."Content", '')),'B'));
    END IF;
    RETURN NULL;
END;
$Content$ LANGUAGE plpgsql;`
	dbGoodsFunc := `CREATE OR REPLACE FUNCTION process_goods_full_text_search() RETURNS TRIGGER AS
$Goods$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        Delete from "FullTextSearch" where "TID" = OLD."ID" and "Type" = 'product';
    ELSIF (TG_OP = 'UPDATE') THEN
        update "FullTextSearch"
        set "UpdatedAt"=NEW."UpdatedAt",
            "Title"=NEW."Title",
            "Content"=NEW."Introduce",
            "Picture"=json_array_element(NEW."Images",0),
            "Index"=setweight(to_tsvector('english', coalesce(NEW."Title", '')),'A') || setweight(to_tsvector('english', coalesce(NEW."Introduce", '')),'B')
        where "TID" = NEW."ID"
          and "Type" = 'product';
    ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO "FullTextSearch" ("ID", "CreatedAt", "UpdatedAt", "OID", "TID", "Title", "Content", "Picture",
                                      "Type", "Index")
        values (DEFAULT, NEW."CreatedAt", NEW."UpdatedAt", NEW."OID", NEW."ID", NEW."Title", NEW."Introduce",
                json_array_element(NEW."Images",0), 'product',
                setweight(to_tsvector('english', coalesce(NEW."Title", '')),'A') || setweight(to_tsvector('english', coalesce(NEW."Introduce", '')),'B'));
    END IF;
    RETURN NULL;
END;
$Goods$ LANGUAGE plpgsql;`

	dbAddContentFunc := `CREATE TRIGGER "Content" AFTER INSERT OR UPDATE OR DELETE ON "Content" FOR EACH ROW EXECUTE FUNCTION process_content_full_text_search();`
	dbAddGoodsFunc := `CREATE TRIGGER "Goods" AFTER INSERT OR UPDATE OR DELETE ON "Goods" FOR EACH ROW EXECUTE FUNCTION process_goods_full_text_search();`

	for index := range models {

		if _database.Migrator().HasTable(models[index]) == false {
			//_database.Set("gorm:table_options", "AUTO_INCREMENT=1000").CreateTable(models[index])
			if err := _database.Migrator().CreateTable(models[index]); err != nil {
				panic(err)
			}
			if models[index].TableName() == (model.FullTextSearch{}).TableName() {
				if err = _database.Exec(`create index idx_FullTextSearch_Index on "FullTextSearch" using gin("Index")`).Error; err != nil {
					panic(err)
				}
				if err = _database.Exec(dbContentFunc).Error; err != nil {
					panic(err)
				}
				if err = _database.Exec(dbGoodsFunc).Error; err != nil {
					panic(err)
				}
				if err = _database.Exec(dbAddContentFunc).Error; err != nil {
					panic(err)
				}
				if err = _database.Exec(dbAddGoodsFunc).Error; err != nil {
					panic(err)
				}
			}
		}
		if !environments.Release() {
			glog.Debug("migrate:", models[index].TableName())
			if err := _database.AutoMigrate(models[index]); err != nil {
				panic(err)
			}
			if err = _database.Exec(dbContentFunc).Error; err != nil {
				panic(err)
			}
			if err = _database.Exec(dbGoodsFunc).Error; err != nil {
				panic(err)
			}
		}

	}

	var _manager model.Manager
	_database.Where(&model.Manager{Account: "manager"}).First(&_manager)
	if _manager.ID == 0 {
		a := model.Manager{Account: "manager", PassWord: encryption.Md5ByString("274455411")}
		if err = _database.Create(&a).Error; err != nil {
			return err
		}
	}

	contentTypeList := []model.ContentType{
		{Type: model.ContentTypeContents, Label: "文章列表"},
		{Type: model.ContentTypeContent, Label: "独立文章"},
		{Type: model.ContentTypeIndex, Label: "首页"},
		{Type: model.ContentTypeGallery, Label: "画廊"},
		{Type: model.ContentTypeProducts, Label: "产品"},
		{Type: model.ContentTypeBlog, Label: "博客"},
		{Type: model.ContentTypePage, Label: "页面"},
	}
	for index := range contentTypeList {
		var _contenttype = contentTypeList[index]
		_database.Where(&model.ContentType{Type: _contenttype.Type}).First(&_contenttype)
		if _contenttype.ID == 0 {
			if err = _database.Create(&_contenttype).Error; err != nil {
				return err
			}
		}
	}

	if !environments.Release() {
		go func() {
			//rebuildFullTextSearch()
		}()
	}
	cache.Init()
	return nil
}
func rebuildFullTextSearch() {
	var err error
	var goodsList []model.Goods
	singleton.Orm().Model(model.Goods{}).Find(&goodsList)
	for _, v := range goodsList {
		var picture string
		if len(v.Images) > 0 {
			picture = v.Images[0]
		}

		fts := model.FullTextSearch{}
		singleton.Orm().Model(model.FullTextSearch{}).Where(`"TID"=? and "Type"=?`, v.ID, model.FullTextSearchTypeProducts).First(&fts)

		fts.CreatedAt = v.CreatedAt
		fts.UpdatedAt = v.UpdatedAt
		fts.OID = v.OID
		fts.TID = v.ID
		fts.Title = v.Title
		fts.Content = util.TrimHtml(v.Introduce)
		fts.Picture = picture
		fts.Type = model.FullTextSearchTypeProducts

		if err = singleton.Orm().Model(&fts).Save(&fts).Error; err != nil {
			panic(err)
		}
		if err = singleton.Orm().Exec(fmt.Sprintf(`UPDATE "FullTextSearch" SET "Index" = setweight(to_tsvector('english', coalesce("Title",'')),'A') || setweight(to_tsvector('english', coalesce("Content",'')),'B') WHERE "ID" = '%d'`, fts.ID)).Error; err != nil {
			panic(err)
		}
	}

	var contentList []model.Content
	singleton.Orm().Model(model.Content{}).Find(&contentList)
	for _, v := range contentList {

		fts := model.FullTextSearch{}
		singleton.Orm().Model(model.FullTextSearch{}).Where(`"TID"=? and "Type"=?`, v.ID, model.FullTextSearchTypeContent).First(&fts)

		fts.CreatedAt = v.CreatedAt
		fts.UpdatedAt = v.UpdatedAt
		fts.OID = v.OID
		fts.TID = v.ID
		fts.Title = v.Title
		fts.Content = util.TrimHtml(v.Content)
		fts.Picture = v.Picture
		fts.Type = model.FullTextSearchTypeContent

		if err = singleton.Orm().Model(&fts).Save(&fts).Error; err != nil {
			panic(err)
		}
		if err = singleton.Orm().Exec(fmt.Sprintf(`UPDATE "FullTextSearch" SET "Index" = setweight(to_tsvector('english', coalesce("Title",'')),'A') || setweight(to_tsvector('english', coalesce("Content",'')),'B') WHERE "ID" = '%d'`, fts.ID)).Error; err != nil {
			panic(err)
		}
	}
}
