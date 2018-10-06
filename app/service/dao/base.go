package dao

import (
	"strings"

	"dandelion/app/play"

	"github.com/jinzhu/gorm"
)

type BaseDao struct {
}

/*func (b BaseDao) UnscopedDeleteWhere(DB *gorm.DB, target interface{}, where interface{}, args ...interface{}) error {
	//db.Where("email LIKE ?", "%jinzhu%").Delete(Email{})
	//DB *gorm.DB, target interface{}, where interface{}
	return DB.Unscoped().Where(where, args...).Delete(target).Error
}*/
/*func (b BaseDao) UnscopedDelete(DB *gorm.DB, target interface{}, ID uint64) error {

	return DB.Unscoped().Delete(target, "ID=?", ID).Error
}*/
func (b BaseDao) DeleteWhere(DB *gorm.DB, target interface{}, where interface{}, args ...interface{}) error {
	//db.Where("email LIKE ?", "%jinzhu%").Delete(Email{})
	//DB *gorm.DB, target interface{}, where interface{}
	return DB.Where(where, args...).Delete(target).Error
}
func (b BaseDao) Delete(DB *gorm.DB, target interface{}, ID uint64) error {

	return DB.Delete(target, "ID=?", ID).Error
}
func (b BaseDao) Add(DB *gorm.DB, target interface{}) error {

	return DB.Create(target).Error
}
func (b BaseDao) Save(DB *gorm.DB, target interface{}) error {

	return DB.Save(target).Error
}
func (b BaseDao) ChangeModel(DB *gorm.DB, ID uint64, target interface{}) error {

	return DB.Model(target).Where("ID=?", ID).Updates(target).Error
}
func (b BaseDao) ChangeMap(DB *gorm.DB, ID uint64, model interface{}, params map[string]interface{}) error {

	return DB.Model(model).Where("ID=?", ID).Updates(params).Error
}
func (b BaseDao) Get(DB *gorm.DB, ID uint64, target interface{}) error {
	return DB.Where("ID=?", ID).First(target).Error
}
func (b BaseDao) FindAllByOID(DB *gorm.DB, target interface{}, OID uint64) error {

	return DB.Where("OID=?", OID).Find(target).Error
}
func (b BaseDao) FindAll(DB *gorm.DB, target interface{}) error {

	return DB.Find(target).Error
}
func (b BaseDao) FindWhere(DB *gorm.DB, target interface{}, where interface{}, args ...interface{}) error {

	return DB.Model(target).Where(where, args...).Find(target).Error

}
func (b BaseDao) FindWhereByOID(DB *gorm.DB, target interface{}, OID uint64, where interface{}, args ...interface{}) error {

	return DB.Model(target).Where("OID=?", OID).Where(where, args...).Find(target).Error

}
func (b BaseDao) FindOrderWhere(DB *gorm.DB, Order interface{}, target interface{}, where interface{}, args ...interface{}) error {

	return DB.Model(target).Where(where, args...).Order(Order).Find(target).Error
}
func (b BaseDao) FindOrderWhereLength(DB *gorm.DB, Order interface{}, target interface{}, Length int, where interface{}, args ...interface{}) error {

	return DB.Model(target).Where(where, args...).Order(Order).Limit(Length).Find(target).Error
}
func (b BaseDao) FindWherePaging(DB *gorm.DB, Order interface{}, target interface{}, Index int, where interface{}, args ...interface{}) error {

	db := DB.Model(target).Where(where, args...).Order(Order)
	SelectPaging(Index, db, target)
	return nil
	//return DB.Model(target).Where(where, args...).Order(Order).Limit(Length).Find(target).Error
}
func (b BaseDao) FindSelectWherePaging(DB *gorm.DB, Select string, Order interface{}, target interface{}, Offset int, where interface{}, args ...interface{}) (_Total, _Limit, _Offset int) {

	db := DB.Model(target).Select(Select).Where(where, args...).Order(Order)
	if Offset < 0 {
		Offset = 0
	}
	_Total, _Offset = SelectPagingOffset(Offset, db, target)
	_Limit = play.Paging
	return
	//return DB.Model(target).Where(where, args...).Order(Order).Limit(Length).Find(target).Error
}

type Search struct {
	Value string `json:"value"`
	Regex bool   `json:"regex"`
}
type Columns struct {
	Data       string `json:"data"`
	Name       string `json:"name"`
	Searchable bool   `json:"searchable"`
	Orderable  bool   `json:"orderable"`
	Search     Search `json:"search"`
}
type Order struct {
	Column int    `json:"column"`
	Dir    string `json:"dir"`
}
type Custom struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}
type Datatables struct {
	//Columns []map[string]interface{} `schema:"columns"`
	Columns []Columns `json:"columns"`
	Customs []Custom  `json:"Customs"`
	Order   []Order   `json:"order"`
	Start   int       `json:"start"`
	Length  int       `json:"length"`
	Search  Search    `json:"search"`
	Draw    int       `json:"draw"`
}

/*`{"draw":1,"columns":[
{"data":"ID","name":"","searchable":true,"orderable":true,"search":{"value":"","regex":false}},
{"data":"Name","name":"","searchable":true,"orderable":true,"search":{"value":"","regex":false}},
{"data":"Grade","name":"","searchable":true,"orderable":true,"search":{"value":"","regex":false}},
{"data":"Province","name":"","searchable":true,"orderable":true,"search":{"value":"","regex":false}},
{"data":"City","name":"","searchable":true,"orderable":true,"search":{"value":"","regex":false}},
{"data":"County","name":"","searchable":true,"orderable":true,"search":{"value":"","regex":false}}
],"order":[
{"column":0,"dir":"asc"}
],"start":0,"length":10,"search":{"value":"","regex":false}}`*/
func (b BaseDao) DatatablesListOrder(Orm *gorm.DB, params *Datatables, target interface{}, OID uint64) (draw int, recordsTotal int, recordsFiltered int, list interface{}) {

	//"draw": 1,
	//"recordsTotal": 57,
	//"recordsFiltered": 57,

	draw = params.Draw

	selectFileds := make([]string, 0)
	searchableFileds := make([]string, 0)
	wheres := make([]string, 0)
	//map[string]interface{}{"name": "jinzhu", "age": 20}

	for _, value := range params.Columns {
		if !strings.EqualFold(value.Data, "") {
			selectFileds = append(selectFileds, value.Data)
		}

		if !strings.EqualFold(value.Search.Value, "") {
			wheres = append(wheres, value.Data+"="+value.Search.Value)

		}
		if value.Searchable && !strings.EqualFold(value.Data, "") {
			searchableFileds = append(searchableFileds, value.Data)
		}

	}

	db := Orm.Select(selectFileds)
	for _, value := range params.Order {
		if !strings.EqualFold(params.Columns[value.Column].Data, "") {
			db = db.Order(params.Columns[value.Column].Data + " " + value.Dir)
		}
	}

	if !strings.EqualFold(params.Search.Value, "") {
		db = db.Where("CONCAT("+strings.Join(searchableFileds, ",")+") like ?", "%"+params.Search.Value+"%").Where(strings.Join(wheres, " and "))
	} else {
		db = db.Where(strings.Join(wheres, " and "))
	}
	if OID != 0 {
		db = db.Where("OID=?", OID)
	}

	if len(params.Customs) > 0 {
		for _, value := range params.Customs {
			db = db.Where(value.Name + value.Value)
		}

	}

	db.Limit(params.Length).Offset(params.Start).Find(target).Offset(0).Count(&recordsTotal)

	recordsFiltered = recordsTotal
	list = target
	return
}
