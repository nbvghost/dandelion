package service

import (
	"dandelion/app/play"
	"dandelion/app/service/dao"
	"strconv"
	"time"
)

type TimeTaskService struct {
	dao.BaseDao
	Transfers TransfersService
	Wx        WxService
	Orders    OrdersService
}

func init() {
	//TimeTaskService := TimeTaskService{}
	//go TimeTaskService.QueryTask()
}
func (self TimeTaskService) QueryTask() {

	go func() {

		c := time.Tick(60 * time.Second)
		for range c {
			self.QuerySupplyOrdersTask()
			self.QueryOrdersTask()
		}

	}()

	c := time.Tick(15 * time.Second)
	for range c {
		//fmt.Printf("在线人数：%v\n", len(gweb.Sessions.Data))
		self.QueryTransfersTask()
	}

}
func (self TimeTaskService) QueryOrdersTask() {
	Orm := dao.Orm()
	var ordersList []dao.Orders
	self.Orders.FindWhere(Orm, &ordersList, "Status<>? and Status<>? and Status<>? and Status<>?", play.OS_OrderOk, play.OS_CancelOk, play.OS_Delete, play.OS_Closed)
	for _, value := range ordersList {

		if value.IsPay == 0 {

			//当前状态为没有支付，去检测一下，订单状态。
			su, result := self.Wx.OrderQuery(value.OrderNo)
			if su {
				TotalFee, _ := strconv.ParseUint(result["total_fee"], 10, 64)
				//OrderNo := result["out_trade_no"]
				//TimeEnd := result["time_end"]
				//attach := result["attach"]
				self.Orders.OrderNotify(TotalFee, result["out_trade_no"], result["time_end"], result["attach"])
				continue
			}
		}
		self.Orders.AnalysisOrdersStatus(value.ID)
	}
}
func (self TimeTaskService) QuerySupplyOrdersTask() {
	Orm := dao.Orm()
	var supplyOrdersList []dao.SupplyOrders
	self.Orders.FindWhere(Orm, &supplyOrdersList, "IsPay=?", 0)
	for _, value := range supplyOrdersList {

		su, result := self.Wx.OrderQuery(value.OrderNo)
		if su {

			TotalFee, _ := strconv.ParseUint(result["total_fee"], 10, 64)
			//OrderNo := result["out_trade_no"]
			//TimeEnd := result["time_end"]
			//attach := result["attach"]
			self.Orders.OrderNotify(TotalFee, result["out_trade_no"], result["time_end"], result["attach"])
			//self.Orders.OrderNotify(result)
		}

	}

}

//查询提现状态
func (self TimeTaskService) QueryTransfersTask() {
	Orm := dao.Orm()
	var transfersList []dao.Transfers
	self.Transfers.FindWhere(Orm, &transfersList, "IsPay=?", 0)
	for _, value := range transfersList {
		su := self.Wx.GetTransfersInfo(value)
		if su {
			self.Transfers.ChangeModel(Orm, value.ID, &dao.Transfers{IsPay: 1})
		} else {
			if time.Now().Unix() > value.CreatedAt.Add(30*time.Hour*24).Unix() {
				self.Transfers.ChangeModel(Orm, value.ID, &dao.Transfers{IsPay: 2})
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
					err := Appointment.ChangeModel(service.Orm, v.ID, &dao.Appointment{State: 4})
					glog.Trace(err)
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
					glog.Trace(err)
					err = Appointment.ChangeModel(service.Orm, v.ID, &dao.Appointment{State: 1, PayDate: t})
					glog.Trace(err)
				}
			} else {
				if (time.Now().Unix() - v.CreatedAt.Unix()) > 30*60 {
					err := Appointment.ChangeModel(service.Orm, v.ID, &dao.Appointment{State: 2})
					glog.Trace(err)
				}
			}

		}
		//fmt.Println(list)
	}
}*/
