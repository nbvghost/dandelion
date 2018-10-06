package service

import (
	"dandelion/app/service/dao"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestUserService_SettlementUser(t *testing.T) {
	type args struct {
		UserID    uint64
		Brokerage uint64
		Order  dao.Orders
	}
	order:=OrdersService{}.GetOrdersByID(2000)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "sdfsdfsd", args: args{UserID: 1000, Brokerage: 6000, Order:order}, wantErr: true},
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

func TestUserService_Leve1(t *testing.T) {
	/*service := UserService{}
	UserID:=uint64(1000)
	fmt.Println(service.Leve1(UserID))
	fmt.Println(service.Leve2(UserID))
	fmt.Println(service.Leve3(UserID))
	fmt.Println(service.Leve4(UserID))
	fmt.Println(service.Leve5(UserID))
	fmt.Println(service.Leve6(UserID))*/
}

func TestUserService_FirstSettlementUserBrokerage(t *testing.T) {
	type fields struct {
		BaseDao       dao.BaseDao
		Configuration ConfigurationService
		GiveVoucher   GiveVoucherService
		CardItem      CardItemService
	}
	type args struct {
		Orm       *gorm.DB
		Brokerage uint64
		Order   dao.Orders
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:"ddd",
			fields:fields{},
			args:args{Orm:dao.Orm(),Brokerage:1000,Order:dao.Orders{OID:2000,UserID:1000}},
			wantErr:true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserService{
				BaseDao:       tt.fields.BaseDao,
				Configuration: tt.fields.Configuration,
				GiveVoucher:   tt.fields.GiveVoucher,
				CardItem:      tt.fields.CardItem,
			}
			if err := service.FirstSettlementUserBrokerage(tt.args.Orm, tt.args.Brokerage, tt.args.Order); (err != nil) != tt.wantErr {
				t.Errorf("UserService.FirstSettlementUserBrokerage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
