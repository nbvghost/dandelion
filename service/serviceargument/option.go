package serviceargument

import (
	"fmt"
	"strings"
)

type OptionsType string

func (m OptionsType) String() string {
	return string(m)
}
func NewOptionsType(v string) OptionsType {
	return OptionsType(v)
}

const (
	OptionsTypeAttribute     OptionsType = "att"
	OptionsTypeSpecification OptionsType = "spe"
	OptionsTypePackageNum    OptionsType = "pac"
	OptionsTypeWeight        OptionsType = "wei"
	OptionsTypePrice         OptionsType = "pri"
)

type OptionValue struct {
	Value string
	Count int
}

func (m OptionValue) GetKey(o Option) string {
	return fmt.Sprintf("%s-%s-%s", o.Type, o.Key, m.Value)
}

type Option struct {
	Type  OptionsType
	Key   string
	Label string
	Value []OptionValue
}

type Options struct {
	Attributes []Option
}

func (m *Options) AddAttributes(optionsType OptionsType, key string, label string, value string) {
	var has bool
	for i := 0; i < len(m.Attributes); i++ {
		item := m.Attributes[i]
		if strings.EqualFold(string(item.Type), string(optionsType)) && strings.EqualFold(string(item.Key), key) {
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
				m.Attributes[i].Value = append(m.Attributes[i].Value, OptionValue{Value: value, Count: 1})
			}
			has = true
			break
		}
	}
	if !has {
		m.Attributes = append(m.Attributes, Option{
			Type:  optionsType,
			Label: label,
			Key:   key,
			Value: []OptionValue{{Value: value, Count: 1}},
		})
	}
}
