package extends

import (
	"gorm.io/gorm/clause"
	"strings"
)

type OrderMethod string

const (
	OrderMethodASC  OrderMethod = "ASC"
	OrderMethodDESC OrderMethod = "DESC"
)

type Order struct {
	ColumnName string
	Method     OrderMethod
}

func (m *Order) OrderByColumn(defaultField string, defaultDesc bool) clause.OrderByColumn {
	var desc = false
	if m.SortMethod() == OrderMethodDESC {
		desc = true
	}
	if strings.EqualFold(m.SortField(), "") {
		return clause.OrderByColumn{Column: clause.Column{Name: defaultField}, Desc: defaultDesc}
	}
	return clause.OrderByColumn{Column: clause.Column{Name: m.SortField()}, Desc: desc}
}
func (m *Order) SortField() string {
	return m.ColumnName
}
func (m *Order) SortMethod() OrderMethod {
	//'ASC' | 'DESC' | '' | 'ascending' | 'descending'
	if strings.EqualFold(string(m.Method), "ascending") {
		return OrderMethodASC
	} else if strings.EqualFold(string(m.Method), "descending") {
		return OrderMethodDESC
	} else if !strings.EqualFold(string(m.Method), "ASC") && !strings.EqualFold(string(m.Method), "DESC") {
		return OrderMethodASC
	} else {
		return m.Method
	}
}
