package repository

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/gpa"
)

var User = gpa.Bind(&UserRepository{}, &model.User{}).(*UserRepository)

type UserRepository struct {
	gpa.IRepository
	GetByEmail func(email string) (*model.User, error) `gpa:"AutoCrate"`
	//UpdateByAge    func(age int, update *params.Update) *result.Update                                   `gpa:"AutoCreate"`

	GetByPhone func(tel string) (*model.User, error) `gpa:"AutoCreate"`
}

func (u *UserRepository) Repository() gpa.IRepository {
	return u.IRepository
}
