package model

import (
	"errors"
	"runtime/debug"
	"time"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/base"
	"github.com/nbvghost/gpa/types"
)

//限时抢购
type TimeSell struct {
	base.BaseModel
	OID       types.PrimaryKey `gorm:"column:OID"`
	Hash      string           `gorm:"column:Hash;unique"` //同一个Hash表示同一个活动
	BuyNum    int              `gorm:"column:BuyNum"`
	Enable    bool             `gorm:"column:Enable"`
	DayNum    int              `gorm:"column:DayNum"`
	Discount  int              `gorm:"column:Discount"`
	TotalNum  int              `gorm:"column:TotalNum"`
	StartTime time.Time        `gorm:"column:StartTime"`
	StartH    int              `gorm:"column:StartH"`
	StartM    int              `gorm:"column:StartM"`
	EndH      int              `gorm:"column:EndH"`
	EndM      int              `gorm:"column:EndM"`
	//GoodsID   uint    `gorm:"column:GoodsID"`
}

func (ts *TimeSell) BeforeCreate(scope *gorm.DB) (err error) {
	if ts.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(ts.TableName() + ":OID不能为空"))

	}
	return nil
}

//是满足所有的限时抢购的条件
func (ts *TimeSell) IsEnable() bool {
	if ts.ID == 0 {
		return false
	}
	if ts.Enable {
		//时间是否到了
		_beginTime := time.Date(ts.StartTime.Year(), ts.StartTime.Month(), ts.StartTime.Day(), ts.StartH, ts.StartM, 0, 0, ts.StartTime.Location())
		_endTime := time.Date(_beginTime.Year(), _beginTime.Month(), _beginTime.Day(), ts.EndH, ts.EndM, 0, 0, _beginTime.Location()).Add(time.Hour * time.Duration(ts.DayNum*24))
		//_beginTime.Add(time.Hour*time.Duration(ts.DayNum*24))

		if time.Now().Unix() >= _beginTime.Unix() && time.Now().Unix() < _endTime.Unix() {
			nowDate := time.Now()
			_startTime := time.Date(nowDate.Year(), nowDate.Month(), nowDate.Day(), ts.StartH, ts.StartM, 0, 0, nowDate.Location())
			_overTime := time.Date(nowDate.Year(), nowDate.Month(), nowDate.Day(), ts.EndH, ts.EndM, 0, 0, nowDate.Location())

			if time.Now().Unix() >= _startTime.Unix() && time.Now().Unix() < _overTime.Unix() {
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	} else {
		return false
	}
}
func (TimeSell) TableName() string {
	return "TimeSell"
}
