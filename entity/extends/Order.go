package extends

type OrderMethod string

const OrderMethodASC = "ASC"
const OrderMethodDESC = "DESC"

type Order struct {
	ColumnName string
	Method     OrderMethod
}
