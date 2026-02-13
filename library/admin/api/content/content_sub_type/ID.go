package content_sub_type

import (
	"errors"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service"
)

type ID struct {
	Organization *model.Organization `mapping:""`
	Delete       struct {
		ID dao.PrimaryKey `uri:"ID"`
	} `method:"Delete"`
	Put struct {
		ID dao.PrimaryKey `uri:"ID"`
		model.ContentSubType
	} `method:"Put"`
	Get struct {
		ID dao.PrimaryKey `uri:"ID"`
	} `method:"Get"`
}

func (m *ID) HandleDelete(ctx constrain.IContext) (r constrain.IResult, err error) {
	Orm := db.GetDB(ctx)
	//ID := object.ParseUint(context.PathParams["ID"])
	css := repository.ContentSubTypeDao.FindContentSubTypesByParentContentSubTypeID(ctx, m.Delete.ID)
	if len(css) > 0 {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("包含子项内容，无法删除"), "删除成功", nil)}, nil
	}
	articles := repository.ContentDao.FindContentByContentSubTypeID(ctx, m.Delete.ID)
	if len(articles) > 0 {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("包含文章，无法删除"), "删除成功", nil)}, nil
	}

	//item := &model.ContentSubType{}
	err = dao.DeleteByPrimaryKey(Orm, entity.ContentSubType, dao.PrimaryKey(m.Delete.ID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "删除成功", nil)}, err
}

func (m *ID) HandlePut(ctx constrain.IContext) (r constrain.IResult, err error) {
	//Orm := db.GetDB(ctx)
	//ID := object.ParseUint(context.PathParams["ID"])
	//item := &model.ContentSubType{}
	//err = util.RequestBodyToJSON(context.Request.Body, item)
	//if err != nil {
	//	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	//}

	err = service.Content.SaveContentSubType(ctx, m.Organization.ID, &(m.Put.ContentSubType))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}

func (m *ID) Handle(ctx constrain.IContext) (r constrain.IResult, err error) {
	//ContentSubTypeID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)

	Orm := db.GetDB(ctx)
	var menus model.ContentSubType
	var pmenus model.ContentSubType

	Orm.Where(`"ID"=?`, m.Get.ID).First(&menus)

	if menus.ID > 0 {
		Orm.Where(`"ID"=?`, menus.ParentContentSubTypeID).First(&pmenus)
	}
	results := make(map[string]interface{})
	results["ContentSubType"] = menus
	results["ParentContentSubType"] = pmenus

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", results)}, nil
}
