package job

import (
	"context"
	"log"

	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service/internal/wechat"
)

type QueryTransfersTask struct {
	WxService wechat.WxService
	Ctx       context.Context
}

func (m *QueryTransfersTask) Run() error {
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
func (m *QueryTransfersTask) work(wechatConfig *model.WechatConfig) error {
	//Orm := singleton.Orm()
	//var transfersList []model.Transfers
	//m.Transfers.FindWhere(Orm, &transfersList, `"IsPay"=?`, 0)
	transfersList := dao.Find(db.GetDB(m.Ctx), &model.Transfers{}).Where(`"OID"=?`, wechatConfig.OID).Where(`"IsPay"=?`, 0).List()
	for i := range transfersList {
		value := transfersList[i].(*model.Transfers)
		transferBatchGet, err := m.WxService.GetTransfersInfo(value, wechatConfig)
		if err != nil {
			log.Println(err)
			continue
		} else {
			isPay := 0
			if *transferBatchGet.BatchStatus == "FINISHED" {
				isPay = 1
				//WAIT_PAY: 待付款确认。需要付款出资商户在商家助手小程序或服务商助手小程序进行付款确认
				//ACCEPTED:已受理。批次已受理成功，若发起批量转账的30分钟后，转账批次单仍处于该状态，可能原因是商户账户余额不足等。商户可查询账户资金流水，若该笔转账批次单的扣款已经发生，则表示批次已经进入转账中，请再次查单确认
				//PROCESSING:转账中。已开始处理批次内的转账明细单
				//FINISHED:已完成。批次内的所有转账明细单都已处理完成
				//CLOSED:已关闭。可查询具体的批次关闭原因确认
			}
			if *transferBatchGet.BatchStatus == "CLOSED" {
				isPay = 2
			}
			err = dao.UpdateByPrimaryKey(db.GetDB(m.Ctx), entity.Transfers, value.ID, &model.Transfers{IsPay: uint(isPay), Status: *transferBatchGet.BatchStatus})
			if err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}

func NewQueryTransfersTask(context context.Context) Job {
	return &QueryTransfersTask{}
}
