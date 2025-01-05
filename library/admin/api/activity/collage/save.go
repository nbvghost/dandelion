package collage

import (
	"errors"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"

	"github.com/nbvghost/tool"
)

type Save struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		CollageJSON string `form:"Collage"`
	} `method:"Post"`
}

func (m *Save) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Save) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)

	//context.Request.ParseForm()

	//CollageJson := context.Request.FormValue("Collage")
	//GoodsListJson := context.Request.FormValue("GoodsListIDs")

	//GoodsListIDs := make([]uint, 0)
	//err = util.JSONToStruct(GoodsListJson, &GoodsListIDs)
	//if err != nil {
	//	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}
	//}

	tx := db.Orm().Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	//item := &model.Collage{}
	item, err := util.JSONToStruct[*model.Collage](m.Post.CollageJSON)
	if err != nil {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	}
	Hash := tool.UUID()
	//if strings.EqualFold(content_item.Hash, "") {
	if item.ID == 0 {
		//新添加
		/*for _, value := range GoodsListIDs {
			isHaveTS := service.GetCollageByGoodsID(value)
			if isHaveTS.ID != 0 && isHaveTS.OID != company.ID {
				continue
			}

			content_item := &model.Collage{}
			err = util.JSONToStruct(CollageJson, content_item)
			content_item.GoodsID = value
			content_item.Hash = Hash
			content_item.OID = company.ID
			err = service.Save(tx, content_item)

		}*/

		//item := &model.Collage{}
		item, err = util.JSONToStruct[*model.Collage](m.Post.CollageJSON)
		//content_item.GoodsID = value
		item.Hash = Hash
		item.OID = m.Organization.ID
		err = dao.Save(tx, item)
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "提交成功", item)}, err

	} else {
		//修改
		_item := service.Activity.Collage.GetCollageByHash(item.Hash, m.Organization.ID)
		if _item.ID == 0 {
			return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("无法修改"), "", nil)}, err
		}
		_item.Num = item.Num
		_item.Discount = item.Discount
		_item.TotalNum = item.TotalNum
		err = dao.Save(tx, _item)

		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "提交成功", _item)}, err
		/*for _, value := range GoodsListIDs {
			isHaveTS := service.GetCollageByGoodsID(value)
			if isHaveTS.ID != 0 {
				if strings.EqualFold(content_item.Hash, isHaveTS.Hash) && isHaveTS.OID == company.ID {
					_item := &model.Collage{}
					err = util.JSONToStruct(CollageJson, _item)
					_item.GoodsID = value
					_item.Hash = content_item.Hash
					_item.OID = company.ID
					_item.ID = isHaveTS.ID
					err = service.Save(tx, _item)
				}
				continue
			}

			_item := &model.Collage{}
			err = util.JSONToStruct(CollageJson, _item)
			_item.GoodsID = value
			_item.Hash = content_item.Hash
			_item.OID = company.ID
			_item.ID = 0
			err = service.Save(tx, _item)

		}*/

	}

	//return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "提交成功", nil)}
}
