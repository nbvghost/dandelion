package wechat

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/order"
	"github.com/nbvghost/dandelion/service/user"

	"reflect"
	"testing"
)

func TestWxService_OrderDeliveryNotify(t *testing.T) {

	Order := order.OrdersService{}.GetOrdersByID(2003)

	tests := []struct {
		name  string
		Order model.Orders
		want  *result.ActionResult
	}{
		{name: "dsf", Order: Order},
	}
	entity := WxService{}
	orderService := order.OrdersService{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ogs, err := orderService.FindOrdersGoodsByOrdersID(singleton.Orm(), tt.Order.ID)
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
		slUser   model.User
		itemName string
		timeText string
		typeText string
	}
	var _user model.User
	user.UserService{}.Get(singleton.Orm(), 1000, &_user)
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
		Order model.Orders
	}

	var Order model.Orders
	order.OrdersService{}.Get(singleton.Orm(), 2003, &Order)

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

			ogs, err := orderService.FindOrdersGoodsByOrdersID(singleton.Orm(), Order.ID)
			if err != nil {
				t.Errorf("OrdersService.FindOrdersGoodsByOrdersID(%v) = %v", Order.ID, err)
			}
			if got := entity.NewOrderNotify(tt.args.Order, ogs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WxService.NewOrderNotify() = %v, want %v", got, tt.want)
			}
		})
	}
}
