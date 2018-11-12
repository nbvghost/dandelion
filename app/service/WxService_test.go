package service

import (
	"dandelion/app/service/dao"
	"reflect"
	"testing"
)

func TestWxService_OrderDeliveryNotify(t *testing.T) {

	Order := OrdersService{}.GetOrdersByID(2003)

	tests := []struct {
		name  string
		Order dao.Orders
		want  *dao.ActionStatus
	}{
		{name: "dsf", Order: Order},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := WxService{}
			if got := entity.OrderDeliveryNotify(tt.Order); !reflect.DeepEqual(got, tt.want) {
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
	var user dao.User
	UserService{}.Get(dao.Orm(), 1000, &user)
	tests := []struct {
		name string
		args args
		want *dao.ActionStatus
	}{
		{name: "sdfsd", args: args{slUser: user, itemName: "dsfsdfsd", timeText: "sdfsdfsd", typeText: "sdfdsfsd"}},
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
	OrdersService{}.Get(dao.Orm(), 2003, &Order)

	tests := []struct {
		name string
		args args
		want *dao.ActionStatus
	}{
		{name: "dsfdsfsd", args: args{Order: Order}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := WxService{}
			if got := entity.NewOrderNotify(tt.args.Order); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WxService.NewOrderNotify() = %v, want %v", got, tt.want)
			}
		})
	}
}
