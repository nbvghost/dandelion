package content_item

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/dao"
)

type CustomizeField struct {
	Organization *model.Organization `mapping:""`
	Delete       struct {
		ContentItemID dao.PrimaryKey `uri:"ContentItemID"`
	} `method:"Delete"`
	Put struct {
		ContentItemID dao.PrimaryKey `uri:"ContentItemID"`
		sqltype.CustomizeField
	} `method:"Put"`
	Get struct {
		ContentItemID dao.PrimaryKey `uri:"ContentItemID"`
	} `method:"Get"`
}

func (m *CustomizeField) HandleDelete(context constrain.IContext) (constrain.IResult, error) {
	return nil, nil
}

func (m *CustomizeField) Handle(context constrain.IContext) (constrain.IResult, error) {
	panic("implement me")
}

func (m *CustomizeField) HandlePut(context constrain.IContext) (constrain.IResult, error) {

	return nil, nil
}
func (m *CustomizeField) HandleGet(context constrain.IContext) (constrain.IResult, error) {
	return nil, nil
}
