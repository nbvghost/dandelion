package dao

import (
	"fmt"
	"strings"
)

type WhereCondition struct {
	field string
	op    string
	value any
}
type Where struct {
	w []WhereCondition
}

func (m *Where) LikeLeft(field string, value any) *Where {
	m.w = append(m.w, WhereCondition{field: field, op: "like", value: fmt.Sprintf("%%%v", value)})
	return m
}
func (m *Where) LikeRight(field string, value any) *Where {
	m.w = append(m.w, WhereCondition{field: field, op: "like", value: fmt.Sprintf("%v%%", value)})
	return m
}
func (m *Where) Like(field string, value any) *Where {
	m.w = append(m.w, WhereCondition{field: field, op: "like", value: fmt.Sprintf("%%%v%%", value)})
	return m
}
func (m *Where) Le(field string, value any) *Where {
	m.w = append(m.w, WhereCondition{field: field, op: "<=", value: value})
	return m
}
func (m *Where) Ge(field string, value any) *Where {
	m.w = append(m.w, WhereCondition{field: field, op: ">=", value: value})
	return m
}
func (m *Where) Lt(field string, value any) *Where {
	m.w = append(m.w, WhereCondition{field: field, op: "<", value: value})
	return m
}
func (m *Where) Gt(field string, value any) *Where {
	m.w = append(m.w, WhereCondition{field: field, op: ">", value: value})
	return m
}
func (m *Where) Eq(field string, value any) *Where {
	m.w = append(m.w, WhereCondition{
		field: field,
		op:    "=",
		value: value,
	})
	return m
}
func (m *Where) Where(where string) *Where {
	m.w = append(m.w, WhereCondition{
		field: where,
		op:    "",
	})
	return m
}
func (m *Where) In(field string, value ...any) *Where {
	m.w = append(m.w, WhereCondition{
		field: field,
		op:    "in",
		value: value,
	})
	return m
}
func (m *Where) String() string {
	var wheres []string
	for _, condition := range m.w {
		if condition.op == "" {
			wheres = append(wheres, condition.field)
		} else {
			wheres = append(wheres, fmt.Sprintf(`%s %s '%v'`, condition.field, condition.op, condition.value))
		}
	}
	return strings.Join(wheres, " and ")
}
func NewWhere() *Where {
	return &Where{}
}
