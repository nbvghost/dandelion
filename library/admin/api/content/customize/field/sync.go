package field

import (
	"encoding/json"
	"strings"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/tool/object"
)

type Sync struct {
	Organization *model.Organization `mapping:""`
	Get          struct{}            `method:"get"` //添加
	Post         struct {
		model.CustomizeFieldGroup
	} `method:"post"` //添加
}

func (m *Sync) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	return nil, nil
}

func (m *Sync) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {

	customizeFieldList := dao.Find(db.GetDB(ctx), &model.CustomizeField{}).Where(`"GroupID"=? and "OID"=? and "Type"='BLOCK'`, m.Post.ID, m.Organization.ID).List()
	list := repository.ContentDao.FindContentByFieldGroupID(ctx, m.Organization.ID, m.Post.ID)

	tx := db.GetDB(ctx).Begin()

	for i := range list {
		item := list[i]
		fieldDataList := make([]map[string]any, 0)
		err := json.Unmarshal([]byte(item.FieldData), &fieldDataList)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for ii := 0; ii < len(fieldDataList); ii++ {
			fieldData := fieldDataList[ii]
			Parent, ParentOk := fieldData["Parent"]
			if !ParentOk {
				continue
			}
			ParentMap, ParentOk := Parent.(map[string]any)
			if !ParentOk {
				continue
			}
			Type := object.ParseString(ParentMap["Type"])
			ID := dao.PrimaryKey(object.ParseUint(ParentMap["ID"]))

			if strings.EqualFold(Type, "BLOCK") && ID > 0 {
				for _, entity := range customizeFieldList {
					customizeField := entity.(*model.CustomizeField)
					if customizeField.ID == ID {
						fieldDataList[ii]["Parent"].(map[string]any)["Field"] = customizeField.Field
					}
				}
			}
		}

		fieldDataBytes, _ := json.Marshal(fieldDataList)
		item.FieldData = string(fieldDataBytes)
		err = dao.UpdateByPrimaryKey(tx, &model.Content{}, item.ID, map[string]any{"FieldData": item.FieldData})
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return result.NewSuccess("同步成功"), nil
}
