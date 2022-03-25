package repository

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/gpa"
	"github.com/nbvghost/gpa/params"
	"github.com/nbvghost/gpa/types"
)

var Content = gpa.Bind(&ContentRepository{}, &model.Content{}).(*ContentRepository)

type ContentRepository struct {
	gpa.IRepository
	FindByOIDLimit                 func(OID types.PrimaryKey, pagination *params.Limit) (pageIndex int, pageSize int, total int, list []*model.Content, err error)                `gpa:"AutoCreate"`
	FindByOIDAndContentItemIDLimit func(OID, ContentItemID types.PrimaryKey, pagination *params.Limit) (pageIndex int, pageSize int, total int, list []*model.Content, err error) `gpa:"AutoCreate"`
	//FindByOIDAndContentItemIDAndContentSubTypeIDInLimit func(OID, ContentItemID types.PrimaryKey, ContentSubTypeID []types.PrimaryKey, pagination *params.Limit) (pageIndex int, pageSize int, total int, list []*model.Content, err error) `gpa:"AutoCreate"`
}

func (u *ContentRepository) Repository() gpa.IRepository {
	return u.IRepository
}
