package timesell

import (
	"errors"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"

	"github.com/nbvghost/dandelion/entity/model"

	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"

	"github.com/nbvghost/tool"
)

type Change struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		TimeSellJSON string `uri:"TimeSell"`
	} `method:"Post"`
}

func (m *Change) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Change) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)

	//context.Request.ParseForm()

	//TimeSellJson := context.Request.FormValue("TimeSell")
	//GoodsListJson := context.Request.FormValue("GoodsListIDs")

	//GoodsListIDs := make([]uint, 0)
	//err = util.JSONToStruct(GoodsListJson, &GoodsListIDs)
	//if err != nil {
	//return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}
	//}

	tx := db.Orm().Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	//item := model.TimeSell{}
	item, err := util.JSONToStruct[*model.TimeSell](m.Post.TimeSellJSON)
	if err != nil {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	}
	Hash := tool.UUID()
	if item.ID == 0 {
		//新添加
		/*for _, value := range GoodsListIDs {
			isHaveTS := service.GetTimeSellByGoodsID(value)
			if isHaveTS.ID != 0 && isHaveTS.OID != company.ID {
				continue
			}

			content_item := &model.TimeSell{}
			err = util.JSONToStruct(TimeSellJson, content_item)
			content_item.GoodsID = value
			content_item.Hash = Hash
			content_item.OID = company.ID
			err = service.Save(tx, content_item)

		}*/

		//item := &model.TimeSell{}
		item, err = util.JSONToStruct[*model.TimeSell](m.Post.TimeSellJSON)
		//content_item.GoodsID = value
		item.Hash = Hash
		item.OID = m.Organization.ID
		err = dao.Save(tx, item)
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "提交成功", item)}, err
	} else {
		_item := service.Activity.TimeSell.GetTimeSellByHash(item.Hash, m.Organization.ID)
		if _item.ID == 0 {
			return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("无法修改"), "", nil)}, err
		}
		//_item.Hash = content_item.Hash
		_item.BuyNum = item.BuyNum
		_item.Enable = item.Enable
		_item.DayNum = item.DayNum
		_item.Discount = item.Discount
		_item.TotalNum = item.TotalNum
		_item.StartTime = item.StartTime
		_item.StartH = item.StartH
		_item.StartM = item.StartM
		_item.EndH = item.EndH
		_item.EndM = item.EndM
		err = dao.Save(tx, _item)

		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "提交成功", _item)}, err

		//修改
		/*for _, value := range GoodsListIDs {
			isHaveTS := service.GetTimeSellByGoodsID(value)
			if isHaveTS.ID != 0 {
				if strings.EqualFold(content_item.Hash, isHaveTS.Hash) && isHaveTS.OID == company.ID {
					_item := &model.TimeSell{}
					err = util.JSONToStruct(TimeSellJson, _item)
					_item.GoodsID = value
					_item.Hash = content_item.Hash
					_item.OID = company.ID
					_item.ID = isHaveTS.ID
					err = service.Save(tx, _item)
				}
				continue
			}

			_item := &model.TimeSell{}
			err = util.JSONToStruct(TimeSellJson, _item)
			_item.GoodsID = value
			_item.Hash = content_item.Hash
			_item.OID = company.ID
			_item.ID = 0
			err = service.Save(tx, _item)

		}*/

	}

}
