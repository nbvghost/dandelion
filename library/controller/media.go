package controller

import (
	"github.com/lib/pq"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/samber/lo"
)

type Media struct {
	Admin *entity.SessionMappingData `mapping:""`
	Get   struct {
		TargetID dao.PrimaryKey `form:"TargetID"`
		Target   string         `form:"Target"`
	} `method:"Get"`
	Post struct {
		TargetIDList []dao.PrimaryKey
		Target       string
	} `method:"Post"`
	Put struct {
		ID   dao.PrimaryKey
		Tags pq.StringArray
	} `method:"Put"`
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
func (m *Media) HandlePut(ctx constrain.IContext) (constrain.IResult, error) {
	media := dao.GetBy(db.Orm(), &model.Media{}, map[string]any{"OID": m.Admin.OID, "ID": m.Put.ID}).(*model.Media)
	if media.IsZero() {
		return nil, result.NewNotFound()
	}
	err := dao.UpdateByPrimaryKey(db.Orm(), &model.Media{}, media.ID, map[string]any{"Tags": lo.Union(m.Put.Tags)})
	if err != nil {
		return nil, err
	}
	return result.NewSuccess("添加成功"), nil
}
func (m *Media) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	list := dao.Find(db.Orm(), &model.Media{}).Where(`"OID"=? and "TargetID" in (?) and "Target"=?`, m.Admin.OID, m.Post.TargetIDList, m.Post.Target).List()
	return result.NewData(map[string]any{"List": list}), nil
}
func (m *Media) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	list := dao.Find(db.Orm(), &model.Media{}).Where(`"OID"=? and "TargetID"=? and "Target"=?`, m.Admin.OID, m.Get.TargetID, m.Get.Target).Order(`"CreatedAt" desc`).List()
	return result.NewData(map[string]any{"List": list}), nil
}
