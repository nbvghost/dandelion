package wechat

import (
	"github.com/nbvghost/dandelion/app/result"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/service/order"
	"github.com/nbvghost/dandelion/app/service/user"
	"reflect"
	"testing"
)

func TestWxService_OrderDeliveryNotify(t *testing.T) {

	Order := order.OrdersService{}.GetOrdersByID(2003)

	tests := []struct {
		name  string
		Order dao.Orders
		want  *result.ActionResult
	}{
		{name: "dsf", Order: Order},
	}
	entity := WxService{}
	orderService := order.OrdersService{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ogs, err := orderService.FindOrdersGoodsByOrdersID(dao.Orm(), tt.Order.ID)
			if err != nil {
				t.Errorf("OrdersService.FindOrdersGoodsByOrdersID(%v) = %v", tt.Order.ID, err)
			}
			if got := entity.OrderDeliveryNotify(tt.Order, ogs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WxService.OrderDeliveryNotify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWxService_INComeNotify(t *testing.T) {

	type args struct {
		slUser   dao.User
		itemName string
		timeText string
		typeText string
	}
	var _user dao.User
	user.UserService{}.Get(dao.Orm(), 1000, &_user)
	tests := []struct {
		name string
		args args
		want *result.ActionResult
	}{
		{name: "sdfsd", args: args{slUser: _user, itemName: "dsfsdfsd", timeText: "sdfsdfsd", typeText: "sdfdsfsd"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := WxService{}
			if got := entity.INComeNotify(tt.args.slUser, tt.args.itemName, tt.args.timeText, tt.args.typeText); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WxService.INComeNotify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWxService_NewOrderNotify(t *testing.T) {

	type args struct {
		Order dao.Orders
	}

	var Order dao.Orders
	order.OrdersService{}.Get(dao.Orm(), 2003, &Order)

	tests := []struct {
		name string
		args args
		want *result.ActionResult
	}{
		{name: "dsfdsfsd", args: args{Order: Order}},
	}

	entity := WxService{}
	orderService := order.OrdersService{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ogs, err := orderService.FindOrdersGoodsByOrdersID(dao.Orm(), Order.ID)
			if err != nil {
				t.Errorf("OrdersService.FindOrdersGoodsByOrdersID(%v) = %v", Order.ID, err)
			}
			if got := entity.NewOrderNotify(tt.args.Order, ogs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WxService.NewOrderNotify() = %v, want %v", got, tt.want)
			}
		})
	}
}
