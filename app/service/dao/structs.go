package dao

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

func (etfi ExpressTemplateItem) CalculateExpressPrice(et ExpressTemplate, nmw ExpressTemplateNMW) uint64 {

	if strings.EqualFold(et.Drawee, "BUSINESS") {
		return 0
	} else {

		//g
		if strings.EqualFold(et.Type, "GRAM") {

			if nmw.W <= etfi.N {
				return uint64(etfi.M * 100)
			} else {
				wp := float64(nmw.W-etfi.N) / float64(etfi.AN) * float64(etfi.ANM*float64(100))
				return uint64(etfi.M*100) + uint64(math.Floor(wp+0.5))
			}

		} else {
			//件
			if nmw.N <= etfi.N {
				return uint64(etfi.M * 100)
			} else {
				wp := float64(nmw.N-etfi.N) / float64(etfi.AN) * float64(etfi.ANM*100)
				return uint64(etfi.M*100) + uint64(math.Floor(wp+0.5))
			}

		}

	}
}

//[{"Areas":["上海","江西省","山东省"],"Type":"N","N":1,"$$hashKey":"object:67"},
// {"Areas":["海南省","青海省","陕西省"],"Type":"M","M":3,"$$hashKey":"object:70"},
// {"Areas":["新疆维吾尔自治区","重庆","四川省"],"Type":"NM","N":3,"M":3,"$$hashKey":"object:73"}]
type ExpressTemplateFreeItem struct {
	Areas []string
	Type  string
	N     int
	M     float64 //元
}

//et 快递模板
//nmw 包邮方式
func (etfi ExpressTemplateFreeItem) IsFree(et ExpressTemplate, nmw ExpressTemplateNMW) bool {
	//ITEM  KG
	if strings.EqualFold(et.Drawee, "BUSINESS") {
		return true
	} else {
		//g
		if strings.EqualFold(et.Type, "GRAM") {

			switch etfi.Type {
			case "N":
				if nmw.W < etfi.N {
					return true
				} else {
					return false
				}

			case "M":

				if nmw.M >= int(math.Floor(etfi.M*100+0.5)) {
					return true
				} else {
					return false
				}

			case "NM":
				if nmw.W < etfi.N && nmw.M > int(math.Floor(etfi.M*100+0.5)) {
					return true
				} else {
					return false
				}
			}

		} else {
			switch etfi.Type {
			case "N":
				if nmw.N > etfi.N {
					return true
				} else {
					return false
				}

			case "M":

				if nmw.M >= int(math.Floor(etfi.M*100+0.5)) {
					return true
				} else {
					return false
				}

			case "NM":
				if nmw.N > etfi.N && nmw.M > int(math.Floor(etfi.M*100+0.5)) {
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

//退货信息
type RefundInfo struct {
	ShipName    string //退货快递公司
	ShipNo      string //退货快递编号
	HasGoods    bool   //是否包含商品，true=包含商品，false=只有款
	Reason      string //原因
	RefundPrice uint64 //返回金额
}
type GoodsParams struct {
	Name  string
	Value string
}
