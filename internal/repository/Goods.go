package repository

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/gpa"
	"github.com/nbvghost/gpa/params"
	"github.com/nbvghost/gpa/types"
)

var Goods = gpa.Bind(&GoodsRepository{}, &model.Goods{}).(*GoodsRepository)

type GoodsRepository struct {
	gpa.IRepository
	FindByOIDLimit                                  func(OID types.PrimaryKey, pagination *params.Limit) (pageIndex int, pageSize int, total int, list []*model.Goods, err error)                                `gpa:"AutoCreate"`
	FindByOIDAndGoodsTypeIDLimit                    func(OID, GoodsTypeID types.PrimaryKey, pagination *params.Limit) (pageIndex int, pageSize int, total int, list []*model.Goods, err error)                   `gpa:"AutoCreate"`
	FindByOIDAndGoodsTypeIDAndGoodsTypeChildIDLimit func(OID, GoodsTypeID, GoodsTypeChildID types.PrimaryKey, pagination *params.Limit) (pageIndex int, pageSize int, total int, list []*model.Goods, err error) `gpa:"AutoCreate"`
}

func (u *GoodsRepository) Repository() gpa.IRepository {
	return u.IRepository
}