package job

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service/internal/order"
	"github.com/nbvghost/dandelion/service/internal/wechat"
)

type QuerySupplyOrdersTask struct {
	WxService     wechat.WxService
	OrdersService order.OrdersService
	Ctx           context.Context
}

func (m *QuerySupplyOrdersTask) Run() error {
	list := m.WxService.MiniProgram(db.GetDB(m.Ctx))
	for i := range list {
		item := list[i].(*model.WechatConfig)
		err := m.work(item)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}
func (m *QuerySupplyOrdersTask) work(wxConfig *model.WechatConfig) error {
	//Orm := singleton.Orm()
	//var supplyOrdersList []model.SupplyOrders
	//m.Orders.FindWhere(Orm, &supplyOrdersList, `"IsPay"=?`, 0)
	supplyOrdersList := dao.Find(db.GetDB(m.Ctx), &model.SupplyOrders{}).Where(`"IsPay"=?`, 0).List()
	for i := range supplyOrdersList {
		value := supplyOrdersList[i].(*model.SupplyOrders)
		transaction, err := m.WxService.OrderQuery(context.TODO(), value.OrderNo, wxConfig)
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
			if err != nil {
				log.Println(err)
			}
			_, err = m.OrdersService.OrderPaySuccess(m.Ctx, uint(*transaction.Amount.PayerTotal), *transaction.OutTradeNo, *transaction.TransactionId, payTime, model.OrdersTypeSupply)
			if err != nil {
				log.Println(err)
			}
		}

	}
	return nil
}

func NewQuerySupplyOrdersTask(context context.Context) Job {
	return &QuerySupplyOrdersTask{}
}
