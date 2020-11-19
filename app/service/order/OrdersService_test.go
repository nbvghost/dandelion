package order

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/nbvghost/dandelion/app/service/dao"
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

func TestOrdersService_FindOrdersGoodsByCollageUser(t *testing.T) {

	type args struct {
		CollageNo string
	}
	tests := []struct {
		name string
		args args
		want []dao.User
	}{
		{name: "test", args: struct{ CollageNo string }{CollageNo: "9d262ef3926bc83f41258410239ce5ba"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := OrdersService{}
			if got := service.FindOrdersGoodsByCollageUser(tt.args.CollageNo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrdersService.FindOrdersGoodsByCollageUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrdersService_ListCollageRecord(t *testing.T) {

	type args struct {
		UserID int
		Index  int
	}
	tests := []struct {
		name             string
		args             args
		wantList         interface{}
		wantTotalRecords int
	}{
		{name: "dsfsd", args: struct {
			UserID int
			Index  int
		}{UserID: 1000, Index: 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//service := OrdersService{}
			/* gotList := service.ListCollageRecord(tt.args.UserID, tt.args.Index)
			if !reflect.DeepEqual(gotList, tt.wantList) {
				t.Errorf("OrdersService.ListCollageRecord() gotList = %v, want %v", gotList, tt.wantList)
			} */

		})
	}
}
