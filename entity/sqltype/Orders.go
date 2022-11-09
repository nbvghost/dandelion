package sqltype

////是否支付成功,0=未支付，1，支付成功，2过期

type OrdersPostType int

const OrdersPostTypePost OrdersPostType = 1    //邮寄
const OrdersPostTypeOffline OrdersPostType = 2 //线下使用

//1=邮寄，2=线下使用
