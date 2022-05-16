package model

import (
	"math"
	"reflect"

	"gorm.io/gorm"

	_ "github.com/lib/pq"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/glog"
	//"github.com/go-gorp/gorp"
	//_ "github.com/go-sql-driver/mysql"
)

type Profiling struct {
	QueryID  int     `gorm:"column:Query_ID"`
	Duration float64 `gorm:"column:Duration"`
	Query    string  `gorm:"column:Query"`
}

func Init(dsn string) {

	glog.Param.StandardOut = true
	glog.Trace("DSN", dsn)

	//conf.Config.DBUrl = "root:274455411@tcp(miniapp.nutsy.cc:3306)/dandelion?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local"

	//postgres.Jsonb{}

	//_database.Exec("SET NAMES utf8mb4")
	//_database.Exec("SET GLOBAL GROUP_CONCAT_MAX_LEN=1844674407370954752")
	//_database.Exec("SET SESSION GROUP_CONCAT_MAX_LEN=1844674407370954752")
	//_database.Exec("SET GLOBAL max_allowed_packet=1844674407370954752")
	//err = _database.Exec("set profiling = 1").Error
	//glog.Error(err)
	go func() {

		/*for {
			// Show PROFILES;
			var profilings []Profiling
			_database.Raw("Show PROFILES").Scan(&profilings)
			if len(profilings) > 1 {

				list := &collections.ListInterface{}
				for k := range profilings {
					list.Append(profilings[k])
				}
				list.SortFunc = func(i, j int) bool {
					a := list.Collection[i].(Profiling)
					b := list.Collection[j].(Profiling)

					if a.Duration > b.Duration {
						return true
					} else {
						return false
					}

				}
				sort.Sort(list)
				fmt.Println(list.Collection[0])

			}
			time.Sleep(3 * time.Second)
		}*/

	}()

}

type IDataBaseFace interface {
	TableName() string
}

func Paging(db *gorm.DB, pageIndex, pageSize int, dataType IDataBaseFace) result.Pager {
	if pageIndex < 0 {
		pageIndex = 0
	}
	if pageSize > 30000 {
		pageSize = 30000
	}
	var _total int64 = 0
	db.Count(&_total)

	GetOffset := func() int {

		//_total, pageIndex, pageSize

		x := float64(_total) / float64(pageSize)
		totalPage := int(math.Ceil(x))
		_pageIndex := pageIndex
		if _pageIndex >= totalPage {
			_pageIndex = totalPage - 1
		}

		return _pageIndex
	}

	_pageIndex := GetOffset()

	dt := reflect.TypeOf(dataType)

	var list = reflect.New(reflect.SliceOf(dt))

	db.Limit(pageSize).Offset(pageIndex * pageSize).Find(list.Interface())

	pager := result.Pager{
		Data:   list.Elem().Interface(),
		Total:  int(_total),
		Limit:  pageSize,
		Offset: _pageIndex,
	}
	return pager.Calculation()
}

/*func SelectPaging(Index int, p *gorm.DB, target interface{}) (Total int) {
	p.Limit(play.Paging).Offset(Index * play.Paging).Find(target).Offset(0).Count(&Total)
	return
}*/

/*func SelectPagingOffset(_Offset int, p *gorm.DB, target interface{}) (Total int, Offset int) {
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
}*/
