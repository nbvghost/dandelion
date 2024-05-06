package sqltype

type DiscountTypeName string

const (
	DiscountTypeNameTimeSell    DiscountTypeName = "TimeSell"
	DiscountTypeNameCardItem    DiscountTypeName = "CardItem"
	DiscountTypeNameCollage     DiscountTypeName = "Collage"
	DiscountTypeNameFullCut     DiscountTypeName = "FullCut"
	DiscountTypeNameGiveVoucher DiscountTypeName = "GiveVoucher"
)

// Discount 优惠商品
type Discount struct {
	Name     string           //优惠名目
	Target   string           //优惠的json数据
	TypeName DiscountTypeName //优惠类型
	Discount uint             //折扣，20%，一般情况下，Discount 和 DiscountAmount 只能在一个
	//DiscountAmount uint             //优惠的金额，单位分，一般情况下，Discount 和 DiscountAmount 只能在一个
}

//extends.Discount{Name: "限时抢购", Target: util.StructToJSON(timeSell), TypeName: "TimeSell", Discount: uint(timeSell.Discount)}
