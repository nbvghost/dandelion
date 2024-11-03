package index

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type GoodsChild struct {
	Get struct {
		GoodsTypeID      dao.PrimaryKey `uri:"GoodsTypeID"`
		GoodsTypeChildID dao.PrimaryKey `uri:"GoodsTypeChildID"`
		Index            int            `form:"Index"`
	} `method:"get"`
	User *model.User `mapping:""`
}

func (m *GoodsChild) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//GoodsTypeID, _ := strconv.ParseUint(context.PathParams["GoodsTypeID"], 10, 64)
	//GoodsTypeChildID, _ := strconv.ParseUint(context.PathParams["GoodsTypeChildID"], 10, 64)
	//Index, _ := strconv.Atoi(context.Request.URL.Query().Get("Index"))
	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)

	//GoodsTypeID       uint  `gorm:"column:GoodsTypeID"`
	//GoodsTypeChildID  uint  `gorm:"column:GoodsTypeChildID"`

	params := &serviceargument.ListQueryParam{}
	if m.Get.GoodsTypeChildID == 0 {
		//sqlWhere = fmt.Sprintf(`"GoodsTypeID"=%v`, m.Get.GoodsTypeID)
		params.GoodsTypeID = m.Get.GoodsTypeID
	} else {
		//sqlWhere = fmt.Sprintf(`"GoodsTypeID"=%v and "GoodsTypeChildID"=%v`, m.Get.GoodsTypeID, m.Get.GoodsTypeChildID)
		params.GoodsTypeID = m.Get.GoodsTypeID
		params.GoodsTypeChildID = m.Get.GoodsTypeChildID
	}

	/**
	Index: 0
	List: [{,…}]
	Size: 10
	Total: 1
	*/

	orderBy := &dao.Sort{}

	return result.NewData(service.Goods.Goods.GoodsList(params, m.User.OID, orderBy.OrderByColumn(`"UpdatedAt"`, true), m.Get.Index+1, 10)), nil //&result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: m.GoodsService.GoodsList(`"UpdatedAt" desc`, m.Get.Index, 10, sqlWhere)}}, nil

	/*if GoodsTypeChildID==0{
		results := controller.Goods.ListGoodsByGoodsTypeID(GoodsTypeID)
		return &gweb.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: results}}
	}else{
		results := controller.Goods.ListGoodsChildByGoodsTypeID(GoodsTypeID, GoodsTypeChildID)
		return &gweb.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: results}}
	}*/
}
