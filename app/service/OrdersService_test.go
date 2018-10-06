package service

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func init() {
	//conf.Config.DBUrl = "root:123456@tcp(127.0.0.1:3306)/shop?charset=utf8&parseTime=True&loc=Local"
}

func TestOrdersService_ListOrdersByStatus(t *testing.T) {

	now := time.Now().Add(-24 * time.Hour)

	service := OrdersService{}
	service.Situation(now.Unix()*1000, now.Unix()*1000)

	vs := VoucherService{}
	vs.Situation(now.Unix()*1000, now.Unix()*1000)

	sgs := ScoreGoodsService{}
	sgs.Situation(now.Unix()*1000, now.Unix()*1000)

	us := UserService{}
	us.Situation(now.Unix()*1000, now.Unix()*1000)

	//orders := OrdersService{}
	//
	//_, tr := orders.ListOrders(0, 0, []string{}, 10, 0)
	////fmt.Println(list)
	//fmt.Println(tr)

}
func TestOrdersService_AnalyseOrdersGoodsList(t *testing.T) {

	/*fmt.Println(strconv.FormatFloat(1.4, 'f', 2, 64))
	fmt.Println(strconv.FormatFloat(1.5, 'f', 2, 64))
	fmt.Println(strconv.FormatFloat(1.6, 'f', 2, 64))*/
	/*fmt.Println(math.Ceil(1.1))
	fmt.Println(math.Ceil(1.2))
	fmt.Println(math.Ceil(1.5))
	fmt.Println(math.Ceil(1.6))
	fmt.Println(math.Ceil(1.9))*/

	fmt.Println(math.Floor(1.1 + 0.5))
	fmt.Println(math.Floor(1.2 + 0.5))
	fmt.Println(math.Floor(1.5 + 0.5))
	fmt.Println(math.Floor(1.6 + 0.5))
	fmt.Println(math.Floor(1.9 + 0.5))
	//fmt.Println(int(float64(1.9)))
	/*fmt.Println(t)
	var oos OrdersService
	address := dao.Addressee{ProvinceName: "山东省"}

	list := make([]dao.OrdersGoods, 0)

	var og dao.OrdersGoods
	og.GoodsID = 2004
	og.SpecificationID = 2006
	og.Quantity = 52
	list = append(list, og)

	fmt.Println(oos.AnalyseOrdersGoodsList(address, list))*/

}
