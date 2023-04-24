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

func (m TimeTaskService) QueryTask(wxConfig *model.WechatConfig) {

	go func() {

		c := time.Tick(60 * time.Second)
		for range c {
			m.QuerySupplyOrdersTask(wxConfig)
		}

	}()

	c := time.Tick(15 * time.Second)
	for range c {
		//fmt.Printf("在线人数：%v\n", len(gweb.Sessions.Data))
		m.QueryTransfersTask(wxConfig)
	}

}

func (m TimeTaskService) QuerySupplyOrdersTask(wxConfig *model.WechatConfig) {
	//Orm := singleton.Orm()
	//var supplyOrdersList []model.SupplyOrders
	//m.Orders.FindWhere(Orm, &supplyOrdersList, `"IsPay"=?`, 0)
	supplyOrdersList := dao.Find(singleton.Orm(), &model.SupplyOrders{}).Where(`"IsPay"=?`, 0).List()
	for i := range supplyOrdersList {
		value := supplyOrdersList[i].(*model.SupplyOrders)
		transaction, err := m.Wx.OrderQuery(context.TODO(), value.OrderNo, wxConfig)
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
			m.Orders.OrderNotify(uint(*transaction.Amount.PayerTotal), *transaction.OutTradeNo, payTime, *transaction.Attach)
			//m.Orders.OrderNotify(result)
		}

	}

}

// 查询提现状态
func (m TimeTaskService) QueryTransfersTask(wxConfig *model.WechatConfig) {
	//Orm := singleton.Orm()
	//var transfersList []model.Transfers
	//m.Transfers.FindWhere(Orm, &transfersList, `"IsPay"=?`, 0)
	transfersList := dao.Find(singleton.Orm(), &model.Transfers{}).Where(`"IsPay"=?`, 0).List()
	for i := range transfersList {
		value := transfersList[i].(*model.Transfers)
		su := m.Wx.GetTransfersInfo(value, wxConfig)
		if su {
			dao.UpdateByPrimaryKey(singleton.Orm(), entity.Transfers, value.ID, &model.Transfers{IsPay: 1})
		} else {
			if time.Now().Unix() > value.CreatedAt.Add(30*time.Hour*24).Unix() {
				dao.UpdateByPrimaryKey(singleton.Orm(), entity.Transfers, value.ID, &model.Transfers{IsPay: 2})
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
