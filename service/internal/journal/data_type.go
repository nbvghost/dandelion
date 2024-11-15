package journal

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/util"
)

type dataTypeUser struct {
	UserID dao.PrimaryKey
}

func (m *dataTypeUser) ToJSON() string {
	return util.StructToJSON(m)
}
func (m *dataTypeUser) ToMap() map[string]any {
	return map[string]any{"UserID": m.UserID}
}
func (m *dataTypeUser) GetType() model.UserJournalType {
	return model.UserJournal_Type_USER_LEVE
}
func NewDataTypeUser(UserID dao.PrimaryKey) IDataType {
	return &dataTypeUser{
		UserID: UserID,
	}
}

type dataTypeOrder struct {
	OrdersID dao.PrimaryKey
}

func (m *dataTypeOrder) ToJSON() string {
	return util.StructToJSON(m)
}
func (m *dataTypeOrder) ToMap() map[string]any {
	return map[string]any{"OrdersID": m.OrdersID}
}
func (m *dataTypeOrder) GetType() model.UserJournalType {
	return model.UserJournal_Type_LEVE
}
func NewDataTypeOrder(OrdersID dao.PrimaryKey) IDataType {
	return &dataTypeOrder{
		OrdersID: OrdersID,
	}
}

type dataTypeTransfers struct {
	TransfersOrderNo string
}

func (m *dataTypeTransfers) ToJSON() string {
	return util.StructToJSON(m)
}
func (m *dataTypeTransfers) ToMap() map[string]any {
	return map[string]any{"TransfersOrderNo": m.TransfersOrderNo}
}
func (m *dataTypeTransfers) GetType() model.UserJournalType {
	return model.UserJournal_Type_TX
}
func NewDataTypeTransfers(TransfersOrderNo string) IDataType {
	return &dataTypeTransfers{
		TransfersOrderNo: TransfersOrderNo,
	}
}

type IDataType interface {
	ToJSON() string
	ToMap() map[string]any
	GetType() model.UserJournalType
}
