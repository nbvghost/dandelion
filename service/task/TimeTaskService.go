package task

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/order"
	"github.com/nbvghost/dandelion/service/wechat"
)

type TimeTaskService struct {
	model.BaseDao
	Transfers order.TransfersService
	Wx        wechat.WxService
	Orders    order.OrdersService
}

func init() {
	//TimeTaskService := TimeTaskService{}
	//go TimeTaskService.QueryTask()
}
func (self TimeTaskService) QueryTask(wxConfig *model.WechatConfig) {

	go func() {

		c := time.Tick(60 * time.Second)
		for range c {
			self.QuerySupplyOrdersTask(wxConfig)
			self.QueryOrdersTask(wxConfig)
		}

	}()

	c := time.Tick(15 * time.Second)
	for range c {
		//fmt.Printf("在线人数：%v\n", len(gweb.Sessions.Data))
		self.QueryTransfersTask(wxConfig)
	}

}
func (self TimeTaskService) QueryOrdersTask(wxConfig *model.WechatConfig) {
	Orm := singleton.Orm()
	var ordersList []model.Orders
	self.Orders.FindWhere(Orm, &ordersList, `"Status"<>? and "Status"<>? and "Status"<>? and "Status"<>?`, model.OrdersStatusOrderOk, model.OrdersStatusCancelOk, model.OrdersStatusDelete, model.OrdersStatusClosed)
	for _, value := range ordersList {

		if value.IsPay == 0 {

			//当前状态为没有支付，去检测一下，订单状态。
			transaction, err := self.Wx.OrderQuery(context.TODO(), value.OrderNo, wxConfig)
			if err != nil {
				log.Println(err)
				continue
			}

			if strings.EqualFold(*transaction.TradeState, "SUCCESS") {
				//TotalFee, _ := strconv.ParseUint(result["total_fee"], 10, 64)
				//OrderNo := result["out_trade_no"]
				//TimeEnd := result["time_end"]
				//attach := result["attach"]
				payTime, err := time.ParseInLocation("2006-01-02T15:04:05-07:00", *transaction.SuccessTime, time.Local)
				log.Println(err)
				self.Orders.OrderNotify(uint(*transaction.Amount.PayerTotal), *transaction.OutTradeNo, payTime, *transaction.Attach)
				continue
			}
		}
		err := self.Orders.AnalysisOrdersStatus(value.ID, wxConfig)
		if err != nil {
			log.Println(err)
		}
	}
}
func (self TimeTaskService) QuerySupplyOrdersTask(wxConfig *model.WechatConfig) {
	Orm := singleton.Orm()
	var supplyOrdersList []model.SupplyOrders
	self.Orders.FindWhere(Orm, &supplyOrdersList, `"IsPay"=?`, 0)
	for _, value := range supplyOrdersList {

		transaction, err := self.Wx.OrderQuery(context.TODO(), value.OrderNo, wxConfig)
		if err != nil {
			log.Println(err)
			continue
		}
		if strings.EqualFold(*transaction.TradeState, "SUCCESS") {

			//TotalFee, _ := strconv.ParseUint(result["total_fee"], 10, 64)
			//OrderNo := result["out_trade_no"]
			//TimeEnd := result["time_end"]
			//attach := result["attach"]
			payTime, err := time.ParseInLocation("2006-01-02T15:04:05-07:00", *transaction.SuccessTime, time.Local)
			log.Println(err)
			self.Orders.OrderNotify(uint(*transaction.Amount.PayerTotal), *transaction.OutTradeNo, payTime, *transaction.Attach)
			//self.Orders.OrderNotify(result)
		}

	}

}

//查询提现状态
func (self TimeTaskService) QueryTransfersTask(wxConfig *model.WechatConfig) {
	Orm := singleton.Orm()
	var transfersList []model.Transfers
	self.Transfers.FindWhere(Orm, &transfersList, `"IsPay"=?`, 0)
	for _, value := range transfersList {
		su := self.Wx.GetTransfersInfo(value, wxConfig)
		if su {
			dao.UpdateByPrimaryKey(Orm, entity.Transfers, value.ID, &model.Transfers{IsPay: 1})
		} else {
			if time.Now().Unix() > value.CreatedAt.Add(30*time.Hour*24).Unix() {
				dao.UpdateByPrimaryKey(Orm, entity.Transfers, value.ID, &model.Transfers{IsPay: 2})
			}
		}
	}
}

/*func StartTimeTask() {

	go pay()
	go refund()

}
func refund() {
	Appointment := service.AppointmentService{}
	c := time.Tick(5 * time.Second)
	for range c {
		list := Appointment.AppointmentDao.FindAppointmentByState(service.Orm, 1)
		for _, v := range list {

			if (time.Now().Unix() - v.PayDate.Unix()) > 30*60 {
				suc := wxpay.Refund(v.ID, v.Score)
				if suc {
					err := Appointment.ChangeModel(service.Orm, v.ID, &model.Appointment{State: 4})
					log.Println(err)
				}
			}

		}
	}
}
func pay() {
	Appointment := service.AppointmentService{}
	c := time.Tick(5 * time.Second)
	for range c {

		list := Appointment.AppointmentDao.FindAppointmentByState(service.Orm, 0)
		for _, v := range list {

			return_code, result_code, trade_state, time_end, total_fee := wxpay.OrderQuery(v.ID)

			if strings.EqualFold("SUCCESS", return_code) && strings.EqualFold("SUCCESS", result_code) && strings.EqualFold("SUCCESS", trade_state) {
				if v.Score == total_fee {

					t, err := time.ParseInLocation("20060102150405", time_end, time.Local)
					log.Println(err)
					err = Appointment.ChangeModel(service.Orm, v.ID, &model.Appointment{State: 1, PayDate: t})
					log.Println(err)
				}
			} else {
				if (time.Now().Unix() - v.CreatedAt.Unix()) > 30*60 {
					err := Appointment.ChangeModel(service.Orm, v.ID, &model.Appointment{State: 2})
					log.Println(err)
				}
			}

		}
		//fmt.Println(list)
	}
}*/
