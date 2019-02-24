package dao

import (
	"errors"
	"math"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"dandelion/app/util"

	"github.com/jinzhu/gorm"
)

type WxConfig struct {
	BaseModel
	AppID     string `gorm:"column:AppID"`
	AppSecret string `gorm:"column:AppSecret"`
	//Type           string `gorm:"column:Type"`
	Token          string `gorm:"column:Token"`
	EncodingAESKey string `gorm:"column:EncodingAESKey"`
	MchID          string `gorm:"column:MchID"`
	PayKey         string `gorm:"column:PayKey"`
}

func (WxConfig) TableName() string {
	return "WxConfig"
}

type Manager struct {
	BaseModel
	Account  string `gorm:"column:Account;not null;unique"`
	PassWord string `gorm:"column:PassWord;not null"`
}

func (Manager) TableName() string {
	return "Manager"
}

type Organization struct {
	BaseModel
	Domain       string    `gorm:"column:Domain;not null;unique"`  //三级域名
	Name         string    `gorm:"column:Name;not null"`           //店名
	Amount       uint64    `gorm:"column:Amount;default:'0'"`      //现金
	BlockAmount  uint64    `gorm:"column:BlockAmount;default:'0'"` //冻结现金
	Province     string    `gorm:"column:Province"`                //省
	City         string    `gorm:"column:City"`                    //市
	District     string    `gorm:"column:District"`                //区域
	Address      string    `gorm:"column:Address"`                 //街道地址
	Telephone    string    `gorm:"column:Telephone"`               //手机
	Categories   string    `gorm:"column:Categories"`              //门店的类型
	Longitude    string    `gorm:"column:Longitude"`               //地理位置
	Latitude     string    `gorm:"column:Latitude"`                //地理位置
	Photos       string    `gorm:"column:Photos"`                  //店的图片
	Special      string    `gorm:"column:Special"`                 //特色
	Opentime     string    `gorm:"column:Opentime"`                //营业时间
	Avgprice     int       `gorm:"column:Avgprice"`                //每人平均消费
	Introduction string    `gorm:"column:Introduction"`            //介绍
	Recommend    string    `gorm:"column:Recommend"`               //推荐
	Vip          int       `gorm:"column:Vip;default:'0'"`         //VIP等级
	Expire       time.Time `gorm:"column:Expire"`                  //过期时间
	Link         string    `gorm:"column:Link"`                    //链接
}

func (Organization) TableName() string {
	return "Organization"
}

type Admin struct {
	BaseModel
	OID         uint64    `gorm:"column:OID"`
	Account     string    `gorm:"column:Account;not null;unique"`
	PassWord    string    `gorm:"column:PassWord;not null"`
	Authority   string    `gorm:"column:Authority;default:''"` //json 权限
	LastLoginAt time.Time `gorm:"column:LastLoginAt"`
}

func (Admin) TableName() string {
	return "Admin"
}

type OrdersGoodsInfo struct {
	OrdersGoods OrdersGoods
	Favoured    Favoured
}
type GoodsInfo struct {
	Goods          Goods
	Specifications []Specification
	Favoured       Favoured
}

type Configuration struct {
	BaseModel
	OID uint64 `gorm:"column:OID"`
	K   uint64 `gorm:"column:K"`
	V   string `gorm:"column:V"`
}

func (u *Configuration) BeforeCreate(scope *gorm.Scope) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))
		return nil
	}
	return nil
}
func (Configuration) TableName() string {
	return "Configuration"
}

type UserInfo struct {
	BaseModel
	UserID       uint64    `gorm:"column:UserID"`
	DaySignTime  time.Time `gorm:"column:DaySignTime"`              //最后一次签到
	DaySignCount int       `gorm:"column:DaySignCount;default:'0'"` //签到次数
}

func (UserInfo) TableName() string {
	return "UserInfo"
}

type UserFormIds struct {
	BaseModel
	UserID uint64 `gorm:"column:UserID"` //
	FormId string `gorm:"column:FormId"` //formId 用于发送
}

func (UserFormIds) TableName() string {
	return "UserFormIds"
}

//优惠商品
type Favoured struct {
	Name     string
	Target   string
	TypeName string
	Discount uint64 //折扣，20%
}
type User struct {
	BaseModel
	Name        string    `gorm:"column:Name"`                    //
	OpenID      string    `gorm:"column:OpenID"`                  //
	Tel         string    `gorm:"column:Tel;not null"`            //
	Age         int       `gorm:"column:Age;default:'0'"`         //
	Region      string    `gorm:"column:Region"`                  //
	Amount      uint64    `gorm:"column:Amount;default:'0'"`      //现金
	BlockAmount uint64    `gorm:"column:BlockAmount;default:'0'"` //冻结现金
	Score       uint64    `gorm:"column:Score;default:'0'"`       //积分
	Growth      uint64    `gorm:"column:Growth;default:'0'"`      //成长值
	Portrait    string    `gorm:"column:Portrait"`                //头像
	Gender      int       `gorm:"column:Gender;default:'0'"`      //性别 1男  2女
	Subscribe   int       `gorm:"column:Subscribe;default:'0'"`   //
	LastLoginAt time.Time `gorm:"column:LastLoginAt"`             //
	RankID      uint64    `gorm:"column:RankID"`                  //
	SuperiorID  uint64    `gorm:"column:SuperiorID"`              //
}

func (User) TableName() string {
	return "User"
}

//条件增送优惠卷
type GiveVoucher struct {
	BaseModel
	OID           uint64 `gorm:"column:OID"`
	ScoreMaxValue uint64 `gorm:"column:ScoreMaxValue"`
	VoucherID     uint64 `gorm:"column:VoucherID"`
}

func (u *GiveVoucher) BeforeCreate(scope *gorm.Scope) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))
		return nil
	}
	return nil
}
func (GiveVoucher) TableName() string {
	return "GiveVoucher"
}

//等级
type Rank struct {
	BaseModel
	GrowMaxValue uint64 `gorm:"column:GrowMaxValue"`
	Title        string `gorm:"column:Title"`
	//VoucherID     uint64 `gorm:"column:VoucherID"`
}

func (Rank) TableName() string {
	return "Rank"
}

//购物车
type ShoppingCart struct {
	BaseModel
	UserID        uint64 `gorm:"column:UserID"`
	GSID          string `gorm:"column:GSID"` //GoodsID+""+SpecificationID
	Goods         string `gorm:"column:Goods;type:text"`
	Specification string `gorm:"column:Specification;type:text"`
	Quantity      uint   `gorm:"column:Quantity"` //数量
}

func (ShoppingCart) TableName() string {
	return "ShoppingCart"
}

type Address struct {
	Name         string
	ProvinceName string
	CityName     string
	CountyName   string
	Detail       string
	PostalCode   string
	Tel          string
}

func (addr Address) IsEmpty() bool {

	return strings.EqualFold(addr.Name, "") || strings.EqualFold(addr.Tel, "") || strings.EqualFold(addr.Detail, "")
}

type ExpressTemplateItem struct {
	Areas []string
	N     int
	M     float64 //元
	AN    int
	ANM   float64 //增加，元
}

func (etfi ExpressTemplateItem) CalculateExpressPrice(et ExpressTemplate, nmw ExpressTemplateNMW) uint64 {

	if strings.EqualFold(et.Drawee, "BUSINESS") {
		return 0
	} else {

		//g
		if strings.EqualFold(et.Type, "GRAM") {

			if nmw.W <= etfi.N {
				return uint64(etfi.M * 100)
			} else {
				wp := float64(nmw.W-etfi.N) / float64(etfi.AN) * float64(etfi.ANM*float64(100))
				return uint64(etfi.M*100) + uint64(math.Floor(wp+0.5))
			}

		} else {
			//件
			if nmw.N <= etfi.N {
				return uint64(etfi.M * 100)
			} else {
				wp := float64(nmw.N-etfi.N) / float64(etfi.AN) * float64(etfi.ANM*100)
				return uint64(etfi.M*100) + uint64(math.Floor(wp+0.5))
			}

		}

	}
}

//[{"Areas":["上海","江西省","山东省"],"Type":"N","N":1,"$$hashKey":"object:67"},
// {"Areas":["海南省","青海省","陕西省"],"Type":"M","M":3,"$$hashKey":"object:70"},
// {"Areas":["新疆维吾尔自治区","重庆","四川省"],"Type":"NM","N":3,"M":3,"$$hashKey":"object:73"}]
type ExpressTemplateFreeItem struct {
	Areas []string
	Type  string
	N     int
	M     float64 //元
}

//et 快递模板
//nmw 包邮方式
func (etfi ExpressTemplateFreeItem) IsFree(et ExpressTemplate, nmw ExpressTemplateNMW) bool {
	//ITEM  KG
	if strings.EqualFold(et.Drawee, "BUSINESS") {
		return true
	} else {
		//g
		if strings.EqualFold(et.Type, "GRAM") {

			switch etfi.Type {
			case "N":
				if nmw.W < etfi.N {
					return true
				} else {
					return false
				}

			case "M":

				if nmw.M >= int(math.Floor(etfi.M*100+0.5)) {
					return true
				} else {
					return false
				}

			case "NM":
				if nmw.W < etfi.N && nmw.M > int(math.Floor(etfi.M*100+0.5)) {
					return true
				} else {
					return false
				}
			}

		} else {
			switch etfi.Type {
			case "N":
				if nmw.N > etfi.N {
					return true
				} else {
					return false
				}

			case "M":

				if nmw.M >= int(math.Floor(etfi.M*100+0.5)) {
					return true
				} else {
					return false
				}

			case "NM":
				if nmw.N > etfi.N && nmw.M > int(math.Floor(etfi.M*100+0.5)) {
					return true
				} else {
					return false
				}
			}
		}
	}
	return false
}

type ExpressTemplateNMW struct {
	N int //数量
	M int //金额 分
	W int //重 kG
}
type ExpressTemplateTemplate struct {
	//{"Default":{"Areas":[],"N":4,"M":4,"AN":4,"ANM":4},"Items":[{"Areas":["江西省","上海"],"N":4,"M":4,"AN":4,"ANM":4,"$$hashKey":"object:144"}]}
	Default ExpressTemplateItem
	Items   []ExpressTemplateItem
}
type ExpressTemplate struct {
	BaseModel
	OID      uint64 `gorm:"column:OID"`
	Name     string `gorm:"column:Name"`
	Drawee   string `gorm:"column:Drawee"`             //付款人
	Type     string `gorm:"column:Type"`               //KG  ITEM
	Template string `gorm:"column:Template;type:text"` //json
	Free     string `gorm:"column:Free;type:text"`     //json []
}

func (u *ExpressTemplate) BeforeCreate(scope *gorm.Scope) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))
		return nil
	}
	return nil
}
func (u ExpressTemplate) TableName() string {
	return "ExpressTemplate"
}

//退货信息
type RefundInfo struct {
	ShipName    string //退货快递公司
	ShipNo      string //退货快递编号
	HasGoods    bool   //是否包含商品，true=包含商品，false=只有款
	Reason      string //原因
	RefundPrice uint64 //返回金额
}
type GoodsParams struct {
	Name  string
	Value string
}
type OrdersGoods struct {
	BaseModel
	OID           uint64 `gorm:"column:OID"`
	OrdersGoodsNo string `gorm:"column:OrdersGoodsNo;unique"` //
	Status        string `gorm:"column:Status"`               //OGAskRefund，OGRefundNo，OGRefundOk，OGRefundInfo，OGRefundComplete
	RefundInfo    string `gorm:"column:RefundInfo;type:text"` //RefundInfo json 退款退货信息
	OrdersID      uint64 `gorm:"column:OrdersID"`             //
	//GoodsID         uint64 `gorm:"column:GoodsID"`                     //
	//SpecificationID uint64 `gorm:"column:SpecificationID"`             //
	Goods         string `gorm:"column:Goods;type:text"`         //josn
	Specification string `gorm:"column:Specification;type:text"` //json
	Favoured      string `gorm:"column:Favoured;type:text"`
	//CollageNo     string `gorm:"column:CollageNo"` //拼团码，每个订单都是唯一
	//TimeSellID     uint64 `gorm:"column:TimeSellID"`             //限时抢购ID
	//TimeSell       string `gorm:"column:TimeSell;type:text"` //json
	Quantity       uint   `gorm:"column:Quantity"`       //数量
	CostPrice      uint64 `gorm:"column:CostPrice"`      //单价-原价
	SellPrice      uint64 `gorm:"column:SellPrice"`      //单价-销售价
	TotalBrokerage uint64 `gorm:"column:TotalBrokerage"` //总佣金
	Error          string `gorm:"column:Error"`          //
}

func (og OrdersGoods) AddError(err string) {

	if strings.EqualFold(og.Error, "") {
		og.Error = err
	} else {
		og.Error = og.Error + "|" + err
	}
}
func (OrdersGoods) TableName() string {
	return "OrdersGoods"
}

//充值
type SupplyOrders struct {
	BaseModel
	OID      uint64    `gorm:"column:OID"`
	UserID   uint64    `gorm:"column:UserID"`         //用户ID，支付的用户ID
	OrderNo  string    `gorm:"column:OrderNo;unique"` //订单号
	StoreID  uint64    `gorm:"column:StoreID"`        //目标ID，如果门店充值的话，这个就是门店的ID，如果普通用户充值的话，这个就是用户ID
	Type     string    `gorm:"column:Type"`           //Store/User
	PayMoney uint64    `gorm:"column:PayMoney"`       //支付金额
	IsPay    uint64    `gorm:"column:IsPay"`          //是否支付成功,0=未支付，1，支付成功，2过期
	PayTime  time.Time `gorm:"column:PayTime"`        //支付时间
}

func (u *SupplyOrders) BeforeCreate(scope *gorm.Scope) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))
		return nil
	}
	return nil
}
func (SupplyOrders) TableName() string {
	return "SupplyOrders"
}

//合并支付
type OrdersPackage struct {
	BaseModel
	OrderNo string `gorm:"column:OrderNo;unique"` //订单号
	//OrderList     string `gorm:"column:OrderList;type:text"` //json []
	TotalPayMoney uint64 `gorm:"column:TotalPayMoney"` //支付价
	IsPay         uint64 `gorm:"column:IsPay"`         //是否支付成功,0=未支付，1，支付成功，2过期
	PrepayID      string `gorm:"column:PrepayID"`      //
	UserID        uint64 `gorm:"column:UserID"`        //用户ID
}

func (OrdersPackage) TableName() string {
	return "OrdersPackage"
}

//订单信息
type Orders struct {
	BaseModel
	OID             uint64    `gorm:"column:OID"`             //
	UserID          uint64    `gorm:"column:UserID"`          //用户ID
	PrepayID        string    `gorm:"column:PrepayID"`        //
	IsPay           uint64    `gorm:"column:IsPay"`           //是否支付成功,0=未支付，1，支付成功，2过期
	OrdersPackageNo string    `gorm:"column:OrdersPackageNo"` //订单号
	OrderNo         string    `gorm:"column:OrderNo;unique"`  //订单号
	PayMoney        uint64    `gorm:"column:PayMoney"`        //支付价
	PostType        int       `gorm:"column:PostType"`        //1=邮寄，2=线下使用
	Status          string    `gorm:"column:Status"`          //状态
	ShipNo          string    `gorm:"column:ShipNo"`          //快递单号
	ShipName        string    `gorm:"column:ShipName"`        //快递
	Address         string    `gorm:"column:Address"`         //收货地址 json
	DeliverTime     time.Time `gorm:"column:DeliverTime"`     //发货时间
	ReceiptTime     time.Time `gorm:"column:ReceiptTime"`     //确认收货时间
	RefundTime      time.Time `gorm:"column:RefundTime"`      //申请退款退货时间
	PayTime         time.Time `gorm:"column:PayTime"`         //支付时间
	DiscountMoney   uint      `gorm:"column:DiscountMoney"`   //相关活动的折扣金额，目前只有满减。
	GoodsMoney      uint      `gorm:"column:GoodsMoney"`      //商品总价
	ExpressMoney    uint      `gorm:"column:ExpressMoney"`    //运费
}

func (u *Orders) BeforeCreate(scope *gorm.Scope) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))
		return nil
	}
	return nil
}

func (Orders) TableName() string {
	return "Orders"
}

/*//订单佣金
type OrderBrokerageTemp struct {
	BaseModel
	UserID    uint64 `gorm:"column:UserID"`
	OrderNo   string `gorm:"column:OrderNo"`
	Brokerage uint64 `gorm:"column:Brokerage"`
}

func (OrderBrokerageTemp) TableName() string {
	return "OrderBrokerageTemp"
}*/

//积分兑换产品
type ScoreGoods struct {
	BaseModel
	OID       uint64 `gorm:"column:OID"`
	Name      string `gorm:"column:Name"`
	Score     int    `gorm:"column:Score"`
	Price     uint64 `gorm:"column:Price"`
	Image     string `gorm:"column:Image"`
	Introduce string `gorm:"column:Introduce"`
}

func (u *ScoreGoods) BeforeCreate(scope *gorm.Scope) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))
		return nil
	}
	return nil
}
func (ScoreGoods) TableName() string {
	return "ScoreGoods"
}

type GoodsType struct {
	BaseModel
	//OID  uint64 `gorm:"column:OID"`
	Name string `gorm:"column:Name"`
}

/*func (u *GoodsType) BeforeCreate(scope *gorm.Scope) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))
		return nil
	}
	return nil
}*/
/*func (u *GoodsType) BeforeSave(scope *gorm.Scope) (err error) {
	var gt GoodsType
	scope.DB().Model(u).Where("OID=?", u.OID).Where("Name=?", u.Name).Find(&gt)
	if gt.ID != 0 {
		err = errors.New("名字重复")
	}
	return
}*/
func (GoodsType) TableName() string {
	return "GoodsType"
}

type GoodsTypeChild struct {
	BaseModel
	//OID         uint64 `gorm:"column:OID"`
	Name        string `gorm:"column:Name"`
	Image       string `gorm:"column:Image"`
	GoodsTypeID uint64 `gorm:"column:GoodsTypeID"`
}

/*func (u *GoodsTypeChild) BeforeCreate(scope *gorm.Scope) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))
		return nil
	}
	return nil
}*/
func (GoodsTypeChild) TableName() string {
	return "GoodsTypeChild"
}

//商品规格
type Specification struct {
	BaseModel
	GoodsID     uint64 `gorm:"column:GoodsID"`
	Label       string `gorm:"column:Label"`
	Num         uint64 `gorm:"column:Num"`    //件
	Weight      uint64 `gorm:"column:Weight"` //每件 多少重 g
	Stock       uint   `gorm:"column:Stock"`
	CostPrice   uint64 `gorm:"column:CostPrice"`   //成本价
	MarketPrice uint64 `gorm:"column:MarketPrice"` //市场价
	Brokerage   uint64 `gorm:"column:Brokerage"`   //总佣金
}

func (Specification) TableName() string {
	return "Specification"
}

//商品
type Goods struct {
	BaseModel
	OID              uint64 `gorm:"column:OID"`
	Title            string `gorm:"column:Title"`
	GoodsTypeID      uint64 `gorm:"column:GoodsTypeID"`
	GoodsTypeChildID uint64 `gorm:"column:GoodsTypeChildID"`
	Price            uint64 `gorm:"column:Price"`
	Stock            uint   `gorm:"column:Stock"`
	Hide             uint   `gorm:"column:Hide"`
	Images           string `gorm:"column:Images;type:text;default:'[]'"` //json array
	Videos           string `gorm:"column:Videos;type:text;default:'[]'"` //json array
	Introduce        string `gorm:"column:Introduce;type:text"`
	Pictures         string `gorm:"column:Pictures;type:text;default:'[]'"` //json array
	Params           string `gorm:"column:Params;type:text;default:'[]'"`   //json array
	//TimeSellID        uint64 `gorm:"column:TimeSellID"`                          //
	ExpressTemplateID uint64 `gorm:"column:ExpressTemplateID"` //
	CountSale         uint64 `gorm:"column:CountSale"`         //销售量
}

func (u *Goods) BeforeCreate(scope *gorm.Scope) (err error) {
	var gt Goods
	scope.DB().Model(u).Where("OID=?", u.OID).Where("Title=?", u.Title).Find(&gt)
	if gt.ID != 0 {
		err = errors.New("名字重复")
	}
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))
		return nil
	}
	return
}
func (u Goods) TableName() string {
	return "Goods"
}

//门店库存
type StoreStock struct {
	BaseModel
	StoreID         uint64 `gorm:"column:StoreID"`
	GoodsID         uint64 `gorm:"column:GoodsID"`
	SpecificationID uint64 `gorm:"column:SpecificationID"`
	Stock           uint64 `gorm:"column:Stock"`
	UseStock        uint64 `gorm:"column:UseStock"` //已经使用的量
}

func (u StoreStock) TableName() string {
	return "StoreStock"
}

//门店
type Store struct {
	BaseModel
	OID          uint64  `gorm:"column:OID"`
	Name         string  `gorm:"column:Name"`
	Address      string  `gorm:"column:Address"`
	Latitude     float64 `gorm:"column:Latitude"`
	Longitude    float64 `gorm:"column:Longitude"`
	Phone        string  `gorm:"column:Phone"`
	Amount       uint64  `gorm:"column:Amount;default:'0'"` //现金
	ServicePhone string  `gorm:"column:ServicePhone"`
	OrderPhone   string  `gorm:"column:OrderPhone"`
	ContactName  string  `gorm:"column:ContactName"`
	Introduce    string  `gorm:"column:Introduce"`
	Images       string  `gorm:"column:Images;type:text"`
	Pictures     string  `gorm:"column:Pictures;type:text"`
	Stars        uint64  `gorm:"column:Stars"`      //总星星数量
	StarsCount   uint64  `gorm:"column:StarsCount"` //评分人数
}

func (u *Store) BeforeCreate(scope *gorm.Scope) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))
		return nil
	}
	return nil
}
func (Store) TableName() string {
	return "Store"
}

//限时抢购
type TimeSell struct {
	BaseModel
	OID       uint64    `gorm:"column:OID"`
	Hash      string    `gorm:"column:Hash"` //同一个Hash表示同一个活动
	BuyNum    int       `gorm:"column:BuyNum"`
	Enable    bool      `gorm:"column:Enable"`
	DayNum    int       `gorm:"column:DayNum"`
	Discount  int       `gorm:"column:Discount"`
	TotalNum  int       `gorm:"column:TotalNum"`
	StartTime time.Time `gorm:"column:StartTime"`
	StartH    int       `gorm:"column:StartH"`
	StartM    int       `gorm:"column:StartM"`
	EndH      int       `gorm:"column:EndH"`
	EndM      int       `gorm:"column:EndM"`
	GoodsID   uint64    `gorm:"column:GoodsID"`
}

func (u *TimeSell) BeforeCreate(scope *gorm.Scope) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))
		return nil
	}
	return nil
}

//是满足所有的限时抢购的条件
func (ts TimeSell) IsEnable() bool {

	if ts.GoodsID == 0 {
		return false
	}
	if ts.ID == 0 {
		return false
	}
	if ts.Enable {
		//时间是否到了
		_beginTime := time.Date(ts.StartTime.Year(), ts.StartTime.Month(), ts.StartTime.Day(), ts.StartH, ts.StartM, 0, 0, ts.StartTime.Location())
		_endTime := time.Date(_beginTime.Year(), _beginTime.Month(), _beginTime.Day(), ts.EndH, ts.EndM, 0, 0, _beginTime.Location()).Add(time.Hour * time.Duration(ts.DayNum*24))
		//_beginTime.Add(time.Hour*time.Duration(ts.DayNum*24))

		if time.Now().Unix() >= _beginTime.Unix() && time.Now().Unix() < _endTime.Unix() {
			nowDate := time.Now()
			_startTime := time.Date(nowDate.Year(), nowDate.Month(), nowDate.Day(), ts.StartH, ts.StartM, 0, 0, nowDate.Location())
			_overTime := time.Date(nowDate.Year(), nowDate.Month(), nowDate.Day(), ts.EndH, ts.EndM, 0, 0, nowDate.Location())

			if time.Now().Unix() >= _startTime.Unix() && time.Now().Unix() < _overTime.Unix() {
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	} else {
		return false
	}
}
func (TimeSell) TableName() string {
	return "TimeSell"
}

//拼团记录
type CollageRecord struct {
	BaseModel
	OrderNo       string `gorm:"column:OrderNo"`
	OrdersGoodsNo string `gorm:"column:OrdersGoodsNo"`
	No            string `gorm:"column:No"`
	UserID        uint64 `gorm:"column:UserID"`
	Collager      uint64 `gorm:"column:Collager"`
	//IsPay         uint64 `gorm:"column:IsPay"` //是否支付成功：0=未支付，1=支付成功
}

func (CollageRecord) TableName() string {
	return "CollageRecord"
}

//拼团
type Collage struct {
	BaseModel
	OID      uint64 `gorm:"column:OID"`
	Hash     string `gorm:"column:Hash"`     //同一个Hash表示同一个活动
	Num      int    `gorm:"column:Num"`      //拼团人数
	Discount int    `gorm:"column:Discount"` // 折扣
	TotalNum int    `gorm:"column:TotalNum"` //总拼团产品数量
	GoodsID  uint64 `gorm:"column:GoodsID"`
}

func (u *Collage) BeforeCreate(scope *gorm.Scope) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))
		return nil
	}
	return nil
}
func (Collage) TableName() string {
	return "Collage"
}

//优惠券
type Voucher struct {
	BaseModel
	OID       uint64 `gorm:"column:OID"`
	Name      string `gorm:"column:Name"`
	Amount    uint64 `gorm:"column:Amount"`
	Image     string `gorm:"column:Image"`
	UseDay    int    `gorm:"column:UseDay"`
	Introduce string `gorm:"column:Introduce"`
}

func (u *Voucher) BeforeCreate(scope *gorm.Scope) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))
		return nil
	}
	return nil
}
func (Voucher) TableName() string {
	return "Voucher"
}

//满减
type FullCut struct {
	BaseModel
	OID       uint64 `gorm:"column:OID"`
	Amount    uint64 `gorm:"column:Amount"`
	CutAmount uint64 `gorm:"column:CutAmount"`
}

func (u *FullCut) BeforeCreate(scope *gorm.Scope) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))
		return nil
	}
	return nil
}
func (FullCut) TableName() string {
	return "FullCut"
}

//省市
type District struct {
	BaseModel
	Code string `gorm:"column:Code;primary_key;unique"`
	Name string `gorm:"column:Name"`
}

func (District) TableName() string {
	return "District"
}

//核销记录-user，store
type Verification struct {
	BaseModel
	VerificationNo string `gorm:"column:VerificationNo;unique"` //订单号
	UserID         uint64 `gorm:"column:UserID"`
	Name           string `gorm:"column:Name"`
	Label          string `gorm:"column:Label"`
	CardItemID     uint64 `gorm:"column:CardItemID"`
	StoreID        uint64 `gorm:"column:StoreID"`
	StoreUserID    uint64 `gorm:"column:StoreUserID"` //门店核销管理员的用户ID
	Quantity       uint   `gorm:"column:Quantity"`    //核销数量
}

func (Verification) TableName() string {
	return "Verification"
}

// 卡
type CardItem struct {
	BaseModel
	OrderNo       string    `gorm:"column:OrderNo;unique"` //订单号
	UserID        uint64    `gorm:"column:UserID"`         //
	Type          string    `gorm:"column:Type"`           //OrdersGoods,Voucher,ScoreGoods
	OrdersGoodsID uint64    `gorm:"column:OrdersGoodsID"`  //
	VoucherID     uint64    `gorm:"column:VoucherID"`      //
	ScoreGoodsID  uint64    `gorm:"column:ScoreGoodsID"`   //
	Data          string    `gorm:"column:Data;type:text"` //json数据
	Quantity      uint      `gorm:"column:Quantity"`       //数量
	UseQuantity   uint      `gorm:"column:UseQuantity"`    //已经使用数量
	ExpireTime    time.Time `gorm:"column:ExpireTime"`     //过期时间
	PostType      int       `gorm:"column:PostType"`       //1=邮寄，2=线下使用
}

func (cardItem CardItem) GetNameLabel(DB *gorm.DB) (Name, Label string) {

	switch cardItem.Type {
	case "OrdersGoods":
		var item OrdersGoods
		DB.First(&item, cardItem.OrdersGoodsID)
		var goods Goods
		var specification Specification
		util.JSONToStruct(item.Goods, &goods)
		util.JSONToStruct(item.Specification, &specification)
		Name = goods.Title
		Label = "规格：" + specification.Label + "(" + strconv.FormatFloat(float64(specification.Num)*float64(specification.Weight)/1000, 'f', 2, 64) + "Kg)"
	case "Voucher":
		var item Voucher
		DB.First(&item, cardItem.VoucherID)
		Name = item.Name
		Label = "金额：" + strconv.FormatFloat(float64(item.Amount)/100, 'f', 2, 64) + "元," + "说明：" + item.Introduce

	case "ScoreGoods":
		var item ScoreGoods
		DB.First(&item, cardItem.ScoreGoodsID)
		Name = item.Name
		Label = "积分：" + strconv.FormatUint(uint64(item.Score), 10) + "," + "市场价：" + strconv.FormatFloat(float64(item.Price)/100, 'f', 2, 64) + "元," + "说明：" + item.Introduce
	}

	return Name, Label
}
func (CardItem) TableName() string {
	return "CardItem"
}

type Logger struct {
	BaseModel
	OID   uint64 `gorm:"column:OID"`
	Key   int    `gorm:"column:Key;default:'0'"`
	Title string `gorm:"column:Title"`
	Data  string `gorm:"column:Data"`
}

func (u *Logger) BeforeCreate(scope *gorm.Scope) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))
		return nil
	}
	return nil
}
func (Logger) TableName() string {
	return "Logger"
}

//账目明细
type UserJournal struct {
	BaseModel
	UserID     uint64 `gorm:"column:UserID"`              //受益者
	Name       string `gorm:"column:Name;not null"`       //
	Detail     string `gorm:"column:Detail;not null"`     //
	Type       int    `gorm:"column:Type;default:'0'"`    //ddddd
	Amount     int64  `gorm:"column:Amount;default:'0'"`  //
	Balance    uint64 `gorm:"column:Balance;default:'0'"` //
	FromUserID uint64 `gorm:"column:FromUserID"`          //来源
	DataKV     string `gorm:"column:DataKV;type:text"`    //{Key:"",Value:""}
}

func (UserJournal) TableName() string {
	return "UserJournal"
}

//Organization
//商店账目明细
type OrganizationJournal struct {
	BaseModel
	OID     uint64 `gorm:"column:OID"`                 //OID
	Name    string `gorm:"column:Name;not null"`       //
	Detail  string `gorm:"column:Detail;not null"`     //
	Type    int    `gorm:"column:Type;default:'0'"`    //ddddd
	Amount  int64  `gorm:"column:Amount;default:'0'"`  //
	Balance uint64 `gorm:"column:Balance;default:'0'"` //
	DataKV  string `gorm:"column:DataKV;type:text"`    //{Key:"",Value:""}
}

func (OrganizationJournal) TableName() string {
	return "OrganizationJournal"
}

//账目明细
type StoreJournal struct {
	BaseModel
	Name     string `gorm:"column:Name;not null"`
	Detail   string `gorm:"column:Detail;not null"`
	StoreID  uint64 `gorm:"column:StoreID"`
	Type     int    `gorm:"column:Type;default:'0'"`    //1=自主核销，2=在线充值
	Amount   int64  `gorm:"column:Amount;default:'0'"`  //变动金额
	Balance  uint64 `gorm:"column:Balance;default:'0'"` //变动后的余额
	TargetID uint64 `gorm:"column:TargetID"`
}

func (StoreJournal) TableName() string {
	return "StoreJournal"
}

type Transfers struct {
	BaseModel
	OrderNo    string `gorm:"column:OrderNo;unique"` //订单号
	UserID     uint64 `gorm:"column:UserID"`         //
	StoreID    uint64 `gorm:"column:StoreID"`
	Amount     uint64 `gorm:"column:Amount;default:'0'"` //提现金额
	ReUserName string `gorm:"column:ReUserName"`         //提现用户真实的名字。
	Desc       string `gorm:"column:Desc"`               //提现说明
	IP         string `gorm:"column:IP"`                 //IP
	OpenId     string `gorm:"column:OpenId"`             //OpenId
	IsPay      uint64 `gorm:"column:IsPay"`              //是否支付成功,0=未支付，1，支付成功，2过期
}

func (Transfers) TableName() string {
	return "Transfers"
}

//Score明细
type ScoreJournal struct {
	BaseModel
	Name    string `gorm:"column:Name;not null"`       //
	Detail  string `gorm:"column:Detail;not null"`     //
	UserID  uint64 `gorm:"column:UserID"`              //
	Score   int64  `gorm:"column:Score;default:'0'"`   //变动金额
	Type    int    `gorm:"column:Type;default:'0'"`    //
	Balance uint64 `gorm:"column:Balance;default:'0'"` //变动后的余额
	DataKV  string `gorm:"column:DataKV;type:text"`    //{Key:"",Value:""}
}

func (ScoreJournal) TableName() string {
	return "ScoreJournal"
}

type KV struct {
	Key   string      //
	Value interface{} //
}

//Content   ContentType  ContentSubType

//Menus
type Content struct {
	BaseModel
	OID           uint64 `gorm:"column:OID"`
	Name          string `gorm:"column:Name"`
	Sort          int    `gorm:"column:Sort"`
	ContentTypeID uint64 `gorm:"column:ContentTypeID"`
	Type          string `gorm:"column:Type"`
	Hide          bool   `gorm:"column:Hide"`
}

func (Content) TableName() string {
	return "Content"
}

//MenuType
type ContentType struct {
	BaseModel
	Label string `gorm:"column:Label"`
	Type  string `gorm:"column:Type;unique"`
}

func (ContentType) TableName() string {
	return "ContentType"
}

//Classify
type ContentSubType struct {
	BaseModel
	Name                   string `gorm:"column:Name"`
	ContentID              uint64 `gorm:"column:ContentID"`
	ParentContentSubTypeID uint64 `gorm:"column:ParentContentSubTypeID"`
	Sort                   int    `gorm:"column:Sort"`
}

func (ContentSubType) TableName() string {
	return "ContentSubType"
}

type Article struct {
	BaseModel
	Title            string `gorm:"column:Title"`
	Content          string `gorm:"column:Content;type:text"`
	Introduce        string `gorm:"column:Introduce"`
	Thumbnail        string `gorm:"column:Thumbnail"`
	ContentID        uint64 `gorm:"column:ContentID"`
	ContentSubTypeID uint64 `gorm:"column:ContentSubTypeID"`
	//ContentSubTypeChildID uint64 `gorm:"column:ContentSubTypeChildID"`
	FromUrl string `gorm:"column:FromUrl"`
	Author  string `gorm:"column:Author"`
	Look    int    `gorm:"column:Look"`
}

func (Article) TableName() string {
	return "Article"
}
