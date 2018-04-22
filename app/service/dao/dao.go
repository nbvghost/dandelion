package dao

import (
	"dandelion/app/play"
	"time"

	//"github.com/go-gorp/gorp"
	//_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Pager struct {
	Data   interface{}
	Total  int
	Limit  int
	Offset int
}

//var dbMap *gorp.DbMap

type BaseModel struct {
	ID        uint64    `gorm:"column:ID;primary_key;unique"`
	CreatedAt time.Time `gorm:"column:CreatedAt"`
	UpdatedAt time.Time `gorm:"column:UpdatedAt"`
}

func SelectPaging(Index int, p *gorm.DB, target interface{}, Total *int) {

	//p.Find(target).Count(total).Limit(PagePaging).Offset(pageIndex * PagePaging).Find(target)
	p.Limit(play.Paging).Offset(Index * play.Paging).Find(target).Offset(0).Count(Total)
}
