package mode

import (
	"fmt"
	"github.com/nbvghost/dandelion/library/dao"
	"strings"
)

type OptionsType string

func (m OptionsType) String() string {
	return string(m)
}
func NewOptionsType(v string) OptionsType {
	return OptionsType(strings.ToUpper(v))
}

const (
	OptionsTypeAttribute     OptionsType = "ATTRIBUTE"
	OptionsTypeSpecification OptionsType = "SPECIFICATION"
	OptionsTypePackageNum    OptionsType = "PACKAGE_NUM"
	OptionsTypeWeight        OptionsType = "WEIGHT"
	OptionsTypePrice         OptionsType = "PRICE"
)

type OptionValue struct {
	ID    dao.PrimaryKey
	Value string
	Count int
}

func (m OptionValue) Key(Type OptionsType) string {
	return fmt.Sprintf("%s-%d", Type, m.ID)
}

type Option struct {
	Type  OptionsType
	Key   string
	Value []OptionValue
}

type Options struct {
	Attributes []Option
}

func (m *Options) AddAttributes(optionsType OptionsType, id dao.PrimaryKey, key, value string) {
	var has bool
	for i := 0; i < len(m.Attributes); i++ {
		item := m.Attributes[i]
		if strings.EqualFold(item.Key, key) {
			var hasOptionValue bool
			for ii := range item.Value {
				optionValue := item.Value[ii]
				if strings.EqualFold(optionValue.Value, value) {
					m.Attributes[i].Value[ii].Count = m.Attributes[i].Value[ii].Count + 1
					hasOptionValue = true
					break
				}
			}
			if !hasOptionValue {
				m.Attributes[i].Value = append(m.Attributes[i].Value, OptionValue{ID: id, Value: value, Count: 1})
			}
			has = true
			break
		}
	}
	if !has {
		m.Attributes = append(m.Attributes, Option{
			Type:  optionsType,
			Key:   key,
			Value: []OptionValue{{ID: id, Value: value, Count: 1}},
		})
	}
}
