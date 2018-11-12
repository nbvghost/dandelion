package service

import (
	"testing"

	"dandelion/app/service/dao"
)

func TestUserService_SettlementUser(t *testing.T) {
	type args struct {
		UserID    uint64
		Brokerage uint64
		Order     dao.Orders
	}
	order := OrdersService{}.GetOrdersByID(2000)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "sdfsdfsd", args: args{UserID: 1000, Brokerage: 6000, Order: order}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserService{}
			if err := service.SettlementUser(dao.Orm(), tt.args.Brokerage, tt.args.Order); (err != nil) != tt.wantErr {
				t.Errorf("UserService.SettlementUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
