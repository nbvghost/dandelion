package repository

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/gpa"
)

var Goods = gpa.Bind(&GoodsRepository{}, &model.Goods{}).(*GoodsRepository)

type GoodsRepository struct {
	gpa.IRepository
}

func (u *GoodsRepository) Repository() gpa.IRepository {
	return u.IRepository
}
