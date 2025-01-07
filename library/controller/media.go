package controller

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type Media struct {
	Admin *entity.SessionMappingData `mapping:""`
	Get   struct {
		TargetID dao.PrimaryKey `form:"TargetID"`
		Target   string         `form:"Target"`
	} `method:"Get"`
	Delete struct {
		ID dao.PrimaryKey `form:"ID"`
	} `method:"Delete"`
}

func (m *Media) HandleDelete(ctx constrain.IContext) (constrain.IResult, error) {
	err := dao.DeleteByPrimaryKey(db.Orm(), &model.Media{}, m.Delete.ID)
	if err != nil {
		return nil, err
	}
	return result.NewData(map[string]any{}), nil
}
func (m *Media) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	list := dao.Find(db.Orm(), &model.Media{}).Where(`"OID"=? and "TargetID"=? and "Target"=?`, m.Admin.OID, m.Get.TargetID, m.Get.Target).List()
	return result.NewData(map[string]any{"List": list}), nil
}
