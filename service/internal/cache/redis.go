package cache

import (
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"log"
	"time"
)

func GetCacheContentSubType(ctx constrain.IContext, oid dao.PrimaryKey) []model.ContentSubType {
	key := fmt.Sprintf("db:cache:ContentSubType:%d", oid)

	var contentSubTypeList []model.ContentSubType

	d, _ := ctx.Redis().Get(ctx, key)
	if len(d) == 0 {
		db.Orm().Model(&model.ContentSubType{}).Where(`"OID"=?`, oid).
			Order(`"Sort"`).Order(`"ID"`).
			Find(&contentSubTypeList)
		bytes, err := json.Marshal(&contentSubTypeList)
		if err != nil {
			log.Println(err)
		}
		err = ctx.Redis().Set(ctx, key, string(bytes), time.Hour*12)
		if err != nil {
			log.Println(err)
		}
		return contentSubTypeList
	} else {
		err := json.Unmarshal([]byte(d), &contentSubTypeList)
		if err != nil {
			return contentSubTypeList
		}
		return contentSubTypeList
	}
}

func GetCacheGoodsType(ctx constrain.IContext, oid dao.PrimaryKey) []model.GoodsType {
	key := fmt.Sprintf("db:cache:GoodsType:%d", oid)

	var contentSubTypeList []model.GoodsType

	d, _ := ctx.Redis().Get(ctx, key)
	if len(d) == 0 {
		db.Orm().Model(&model.GoodsType{}).Where(`"OID"=?`, oid).
			Order(`"ID"`).
			Find(&contentSubTypeList)
		bytes, err := json.Marshal(&contentSubTypeList)
		if err != nil {
			log.Println(err)
		}
		err = ctx.Redis().Set(ctx, key, string(bytes), time.Hour*12)
		if err != nil {
			log.Println(err)
		}
		return contentSubTypeList
	} else {
		err := json.Unmarshal([]byte(d), &contentSubTypeList)
		if err != nil {
			return contentSubTypeList
		}
		return contentSubTypeList
	}
}


func GetCacheGoodsTypeChild(ctx constrain.IContext, oid dao.PrimaryKey) []model.GoodsTypeChild {
	key := fmt.Sprintf("db:cache:GoodsTypeChild:%d", oid)

	var contentSubTypeList []model.GoodsTypeChild

	d, _ := ctx.Redis().Get(ctx, key)
	if len(d) == 0 {
		db.Orm().Model(&model.GoodsTypeChild{}).Where(`"OID"=?`, oid).
			Order(`"ID"`).
			Find(&contentSubTypeList)
		bytes, err := json.Marshal(&contentSubTypeList)
		if err != nil {
			log.Println(err)
		}
		err = ctx.Redis().Set(ctx, key, string(bytes), time.Hour*12)
		if err != nil {
			log.Println(err)
		}
		return contentSubTypeList
	} else {
		err := json.Unmarshal([]byte(d), &contentSubTypeList)
		if err != nil {
			return contentSubTypeList
		}
		return contentSubTypeList
	}
}


func GetCacheContentItem(ctx constrain.IContext, oid dao.PrimaryKey) []model.ContentItem {
	key := fmt.Sprintf("db:cache:ContentItem:%d", oid)

	var contentSubTypeList []model.ContentItem

	d, _ := ctx.Redis().Get(ctx, key)
	if len(d) == 0 {
		db.Orm().Model(&model.ContentItem{}).Where(`"OID"=?`, oid).
			Order(`"Sort"`).
			Find(&contentSubTypeList)
		bytes, err := json.Marshal(&contentSubTypeList)
		if err != nil {
			log.Println(err)
		}
		err = ctx.Redis().Set(ctx, key, string(bytes), time.Hour*12)
		if err != nil {
			log.Println(err)
		}
		return contentSubTypeList
	} else {
		err := json.Unmarshal([]byte(d), &contentSubTypeList)
		if err != nil {
			return contentSubTypeList
		}
		return contentSubTypeList
	}
}