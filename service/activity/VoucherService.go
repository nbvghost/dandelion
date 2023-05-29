package activity

import (
	"github.com/nbvghost/dandelion/library/db"
	"time"

	"github.com/nbvghost/dandelion/entity/model"
)

type VoucherService struct {
	model.BaseDao
}

func (service VoucherService) Situation(StartTime, EndTime int64) interface{} {

	st := time.Unix(StartTime/1000, 0)
	st = time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location())
	et := time.Unix(EndTime/1000, 0).Add(24 * time.Hour)
	et = time.Date(et.Year(), et.Month(), et.Day(), 0, 0, 0, 0, et.Location())

	Orm := db.Orm()

	type Result struct {
		TotalMoney uint `gorm:"column:TotalMoney"`
		TotalCount uint `gorm:"column:TotalCount"`
	}

	//select COUNT(ID),SUM(JSON_EXTRACT(Data,'$.Amount')) as TotalMoney from carditem where Type='Voucher' group by Type;

	var result Result

	Orm.Table("CardItem").Select("SUM(JSON_EXTRACT(Data,'$.Amount')) as TotalMoney,COUNT(ID) as TotalCount").Where("CreatedAt>=?", st).Where("CreatedAt<?", et).Where("Type=?", "Voucher").Find(&result)
	//fmt.Println(result)
	return result
}
