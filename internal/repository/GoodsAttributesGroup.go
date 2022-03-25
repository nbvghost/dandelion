package repository

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/gpa"
	"github.com/nbvghost/gpa/types"
)

var GoodsAttributesGroup = gpa.Bind(&GoodsAttributesGroupRepository{}, &model.GoodsAttributesGroup{}).(*GoodsAttributesGroupRepository)

type GoodsAttributesGroupRepository struct {
	gpa.IRepository
	FindByGoodsID       func(goodsID types.PrimaryKey) ([]*model.GoodsAttributesGroup, error)            `gpa:"AutoCrate"`
	GetByGoodsIDAndName func(goodsID types.PrimaryKey, name string) (*model.GoodsAttributesGroup, error) `gpa:"AutoCrate"`
	GetByName           func(name string) (*model.GoodsAttributesGroup, error)                           `gpa:"AutoCrate"`
}

func (u *GoodsAttributesGroupRepository) Repository() gpa.IRepository {

	return u.IRepository
}
