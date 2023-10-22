package model

import (
	"math"
	"strings"
)

type ExpressTemplateItem struct {
	Areas []string
	N     int
	M     float64 //元
	AN    int
	ANM   float64 //增加，元
}

func (m ExpressTemplateItem) CalculateExpressPrice(et *ExpressTemplate, nmw ExpressTemplateNMW) uint {

	if strings.EqualFold(et.Drawee, "BUSINESS") {
		return 0
	} else {

		//g
		if strings.EqualFold(et.Type, "GRAM") {

			if nmw.W <= m.N {
				return uint(m.M * 100)
			} else {
				wp := float64(nmw.W-m.N) / float64(m.AN) * float64(m.ANM*float64(100))
				return uint(m.M*100) + uint(math.Floor(wp+0.5))
			}

		} else {
			//件
			if nmw.N <= m.N {
				return uint(m.M * 100)
			} else {
				wp := float64(nmw.N-m.N) / float64(m.AN) * float64(m.ANM*100)
				return uint(m.M*100) + uint(math.Floor(wp+0.5))
			}

		}

	}
}

// [{"Areas":["上海","江西省","山东省"],"Type":"N","N":1,"$$hashKey":"object:67"},
// {"Areas":["海南省","青海省","陕西省"],"Type":"M","M":3,"$$hashKey":"object:70"},
// {"Areas":["新疆维吾尔自治区","重庆","四川省"],"Type":"NM","N":3,"M":3,"$$hashKey":"object:73"}]
type ExpressTemplateFreeItem struct {
	Areas []string
	Type  string
	N     int
	M     float64 //元
}

// et 快递模板
// nmw 包邮方式
func (m ExpressTemplateFreeItem) IsFree(et *ExpressTemplate, nmw ExpressTemplateNMW) bool {
	//ITEM  KG
	if strings.EqualFold(et.Drawee, "BUSINESS") {
		return true
	} else {
		//g
		if strings.EqualFold(et.Type, "GRAM") {

			switch m.Type {
			case "N":
				if nmw.W < m.N {
					return true
				} else {
					return false
				}

			case "M":

				if nmw.M >= int(math.Floor(m.M*100+0.5)) {
					return true
				} else {
					return false
				}

			case "NM":
				if nmw.W < m.N && nmw.M > int(math.Floor(m.M*100+0.5)) {
					return true
				} else {
					return false
				}
			}

		} else {
			switch m.Type {
			case "N":
				if nmw.N > m.N {
					return true
				} else {
					return false
				}

			case "M":

				if nmw.M >= int(math.Floor(m.M*100+0.5)) {
					return true
				} else {
					return false
				}

			case "NM":
				if nmw.N > m.N && nmw.M > int(math.Floor(m.M*100+0.5)) {
					return true
				} else {
					return false
				}
			}
		}
	}
	return false
}

type ExpressTemplateNMW struct {
	N int //数量
	M int //金额 分
	W int //重 kG
}
type ExpressTemplateTemplate struct {
	//{"Default":{"Areas":[],"N":4,"M":4,"AN":4,"ANM":4},"Items":[{"Areas":["江西省","上海"],"N":4,"M":4,"AN":4,"ANM":4,"$$hashKey":"object:144"}]}
	Default ExpressTemplateItem
	Items   []ExpressTemplateItem
}
