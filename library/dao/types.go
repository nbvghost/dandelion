package dao

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type PrimaryKey uint

func (bm PrimaryKey) String() string {
	return strconv.FormatUint(uint64(bm), 10)
}
func (bm PrimaryKey) IsZero() bool {
	if bm == 0 {
		return true
	} else {
		return false
	}
}
func NewFromString(v string) PrimaryKey {
	n, _ := strconv.ParseUint(v, 10, 64)
	return PrimaryKey(n)
}

// var _ IEntity = (*Entity)(nil)
type IEntity interface {
	TableName() string
	IsZero() bool
	Primary() PrimaryKey
	PrimaryName() string
}

type Entity struct {
	ID        PrimaryKey `gorm:"COMMENT:自增ID;NOT NULL;column:ID;PRIMARY_KEY;AUTOINCREMENT"`                  //条目ID
	CreatedAt time.Time  `gorm:"COMMENT:创建日期;NOT NULL;Type:time;DEFAULT:CURRENT_TIMESTAMP;column:CreatedAt"` //登陆日期
	UpdatedAt time.Time  `gorm:"COMMENT:更新日期;NOT NULL;Type:time;DEFAULT:CURRENT_TIMESTAMP;column:UpdatedAt"` //修改日期,不支持 ON UPDATE:CURRENT_TIMESTAMP 语句
	//DeletedAt *time.Time `gorm:"column:DeletedAt"`             //删除日期
}

/*func (bm Entity) TableName() string {
	//TODO implement me
	panic("implement me")
}*/

func (bm *Entity) Primary() PrimaryKey {
	return bm.ID
}
func (bm *Entity) PrimaryName() string {
	return "ID"
}
func (bm *Entity) IsZero() bool {
	if bm == nil {
		return true
	}
	if bm.ID == 0 && bm.CreatedAt.IsZero() && bm.UpdatedAt.IsZero() {
		return true
	} else {
		return false
	}
}
func NewEntityID(ID PrimaryKey) Entity {
	return Entity{ID: ID}
}

type LikeMod struct {
	InLeft  bool
	InRight bool
}

func NewLikeMod(left bool, right bool) LikeMod {
	return LikeMod{
		InLeft:  left,
		InRight: right,
	}
}

/*func (s LikeMod) Bool() (left bool, right bool) {
	switch s {
	case LikeModNone:
		return false, false
	case LikeModLeft:
		return true, false
	case LikeModRight:
		return false, true
	case LikeModLeftRight:
		return true, true
	}
	return false, false
}*/

//const LikeModNone LikeMod = "NONE"
//const LikeModLeft LikeMod = "LEFT"
//const LikeModRight LikeMod = "RIGHT"
//const LikeModLeftRight LikeMod = "LEFT_RIGHT"

type WhereConditionValue struct {
	IsValue bool //值或点位符,值为=true
	Value   interface{}
}

func NewWhereConditionValue(IsValue bool, Value interface{}) WhereConditionValue {
	return WhereConditionValue{IsValue: IsValue, Value: Value}
}

type OperationType string

const (
	GetOperationType    OperationType = "GET"
	FindOperationType   OperationType = "FIND"
	DeleteOperationType OperationType = "DELETE"
	UpdateOperationType OperationType = "UPDATE"
	InsertOperationType OperationType = "INSERT"
)

func NewOperationType(sql string) (OperationType, error) {
	queryType := strings.ToUpper(sql[:3])
	if queryType != "GET" {
		queryType = strings.ToUpper(sql[:4])
		if queryType != "FIND" {
			queryType = strings.ToUpper(sql[:6])
			if queryType != "DELETE" {
				queryType = strings.ToUpper(sql[:6])
				if queryType != "UPDATE" {
					return "", errors.New("函数必须以[Get,Find,Delete,Update]开头")
				}
			}
		}
		//isMany = true
	}
	return OperationType(queryType), nil
}
