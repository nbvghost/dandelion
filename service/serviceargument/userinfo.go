package serviceargument

import (
	"github.com/nbvghost/dandelion/entity/model"
)

const (
	UserInfoKeyBrokerageLeve1 model.UserInfoKey = "BrokerageLeve1"
	UserInfoKeyBrokerageLeve2 model.UserInfoKey = "BrokerageLeve2"
	UserInfoKeyBrokerageLeve3 model.UserInfoKey = "BrokerageLeve3"
	UserInfoKeyBrokerageLeve4 model.UserInfoKey = "BrokerageLeve4"
	UserInfoKeyBrokerageLeve5 model.UserInfoKey = "BrokerageLeve5"
	UserInfoKeyBrokerageLeve6 model.UserInfoKey = "BrokerageLeve6"
)

type UserInfoKeyStateType string

const (
	UserInfoKeyStateTypeNormal  UserInfoKeyStateType = ""        //
	UserInfoKeyStateTypeClosure UserInfoKeyStateType = "CLOSURE" //封闭
)
