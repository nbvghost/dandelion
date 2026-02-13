package content

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/repository"
)

type SingleGetContentItemIDContentSubTypeID struct {
	Get struct {
		ContentItemID    uint `uri:"ContentItemID"`
		ContentSubTypeID uint `uri:"ContentSubTypeID"`
	} `method:"get"`
}

func (m *SingleGetContentItemIDContentSubTypeID) Handle(ctx constrain.IContext) (r constrain.IResult, err error) {
	//ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	//ContentSubTypeID, _ := strconv.ParseUint(context.PathParams["ContentSubTypeID"], 10, 64)

	article := repository.ContentDao.FindContentByContentItemIDAndContentSubTypeID(ctx, dao.PrimaryKey(m.Get.ContentItemID), dao.PrimaryKey(m.Get.ContentSubTypeID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", article)}, nil
}
