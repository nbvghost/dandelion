package dao

import (
	"dandelion/app/play"
	"time"

	_ "github.com/go-sql-driver/mysql"
	//"github.com/go-gorp/gorp"
	//_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/gweb/tool"
)

var _database *gorm.DB

func Orm() *gorm.DB {

	//return _database.New()
	return _database
}
func init() {
	var err error
	//open:=make(chan bool,1)
	_database, err = gorm.Open("mysql", conf.Config.DBUrl)
	tool.CheckError(err)

	_database.Debug()
	_database.LogMode(true)

	_database.Exec("SET NAMES utf8mb4")
	_database.Exec("SET GLOBAL GROUP_CONCAT_MAX_LEN=1844674407370954752")
	_database.Exec("SET SESSION GROUP_CONCAT_MAX_LEN=1844674407370954752")
	_database.Exec("SET GLOBAL max_allowed_packet=1844674407370954752")

}

type Pager struct {
	Data   interface{}
	Total  int
	Limit  int
	Offset int
}

//var dbMap *gorp.DbMap
type ActionStatus struct {
	Success bool
	Message string
	Data    interface{}
}

func (as *ActionStatus) SmartSuccessData(data interface{}) *ActionStatus {
	as.Message = "SUCCESS"
	as.Success = true
	as.Data = data
	return as
}
func (as *ActionStatus) SmartError(err error, successTxt string, data interface{}) *ActionStatus {

	if err == nil {
		as.Message = successTxt
		as.Success = true
		as.Data = data
	} else {
		as.Message = err.Error()
		as.Success = false
		as.Data = data
	}
	return as
}
func (as *ActionStatus) Smart(success bool, s string, f string) *ActionStatus {
	as.Success = success
	if success {
		as.Message = s
	} else {
		as.Message = f
	}
	return as
}
func (as *ActionStatus) SmartData(success bool, s string, f string, data interface{}) *ActionStatus {
	as.Success = success
	if success {
		as.Message = s
		as.Data = data
	} else {
		as.Message = f
	}
	return as
}

type BaseModel struct {
	ID        uint64     `gorm:"column:ID;primary_key;unique"` //条目ID
	CreatedAt time.Time  `gorm:"column:CreatedAt"`             //登陆日期
	UpdatedAt time.Time  `gorm:"column:UpdatedAt"`             //修改日期
	DeletedAt *time.Time `gorm:"column:DeletedAt"`             //删除日期
	//Delete    int        `gorm:"column:Delete"`                //0=无，1=删除，
}

func SelectPaging(Index int, p *gorm.DB, target interface{}) (Total int) {

	//p.Find(target).Count(total).Limit(PagePaging).Offset(pageIndex * PagePaging).Find(target)
	p.Limit(play.Paging).Offset(Index * play.Paging).Find(target).Offset(0).Count(&Total)
	return
}
func SelectPagingOffset(_Offset int, p *gorm.DB, target interface{}) (Total int, Offset int) {

	//p.Find(target).Count(total).Limit(PagePaging).Offset(pageIndex * PagePaging).Find(target)
	p.Limit(play.Paging).Offset(_Offset).Find(target).Offset(0).Count(&Total)
	if _Offset >= Total {
		Offset = Total
	} else {
		if Total-_Offset < 0 {
			Offset = 0
		} else {
			Offset = _Offset
		}
	}
	return
}
