package dao

import (
	"encoding/json"
	"time"

	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"
)

func OutActionStatus(as *ActionStatus) []byte {

	b, err := json.Marshal(as)
	tool.CheckError(err)
	return b
}
func WriteJSON(context *gweb.Context, as *ActionStatus) {

	b := OutActionStatus(as)

	context.Response.Write(b)
}
func WritePager(context *gweb.Context, pager *Pager) {

	b, err := json.Marshal(pager)
	tool.CheckError(err)
	context.Response.Write(b)
}

/*config := &wechat.Config{
AppID:          "xxxx",
AppSecret:      "xxxx",
Token:          "xxxx",
EncodingAESKey: "xxxx",
Cache:          memCache
}*/

type Manager struct {
	BaseModel
	Account  string `gorm:"column:Account;not null;unique"`
	PassWord string `gorm:"column:PassWord;not null"`
}

func (Manager) TableName() string {
	return "Manager"
}

type TempOrderPack struct {
	CompanyID uint64
	ShopName  string
	Orders    []*Order
}
type OrderPack struct {
	ID        uint64    `gorm:"column:ID;primary_key;unique"`
	CompanyID uint64    `gorm:"column:CompanyID;not null"`
	UserID    uint64    `gorm:"column:UserID;not null"`
	OrderNo   string    `gorm:"column:OrderNo;not null;unique"`
	Position  string    `gorm:"column:Position;default:''"`
	OrderList string    `gorm:"column:OrderList;type:json"` //[Order,Order,Order,Order,Order]
	Tip       string    `gorm:"column:Tip;default:''"`
	PostType  int       `gorm:"column:PostType;default:'0'"` //0=店里支付，1=邮寄
	Address   string    `gorm:"column:Address;default:''"`
	Status    int       `gorm:"column:Status;default:'0'"` //0=订单提交成功 1=支付成功  2=订单结束
	Total     uint64    `gorm:"column:Total;default:'0'"`  //总价
	PayAt     time.Time `gorm:"column:PayAt"`
	PostAt    time.Time `gorm:"column:PostAt"`
	ClosedAt  time.Time `gorm:"column:ClosedAt"`
	CreatedAt time.Time `gorm:"column:CreatedAt"`
}

func (OrderPack) TableName() string {
	return "OrderPack"
}

type Order struct {
	OID          uint64 //原项目ID，类型由Type来决定
	CompanyID    uint64
	UserID       uint64
	Type         int //1=appointment
	Title        string
	CacheContent interface{}
	Params       string //[{Key:'颜色',Value:'黑色'},{Key:'颜色',Value:'黑色'}]
	Count        uint64 //数量
	Price        uint64 //单价
	Total        uint64 //总价
}

type ActionStatus struct {
	Success bool
	Message string
	Data    interface{}
}

func (as *ActionStatus) SmartSuccessData(data interface{}) *ActionStatus {
	as.Message = "SUCCESS"
	as.Success = true
	as.Data = data
	return as
}
func (as *ActionStatus) SmartError(err error, successTxt string, data interface{}) *ActionStatus {

	if err == nil {
		as.Message = successTxt
		as.Success = true
		as.Data = data
	} else {
		as.Message = err.Error()
		as.Success = false
		as.Data = data
	}
	return as
}
func (as *ActionStatus) Smart(success bool, s string, f string) *ActionStatus {
	as.Success = success
	if success {
		as.Message = s
	} else {
		as.Message = f
	}
	return as
}
func (as *ActionStatus) SmartData(success bool, s string, f string, data interface{}) *ActionStatus {
	as.Success = success
	if success {
		as.Message = s
		as.Data = data
	} else {
		as.Message = f
	}
	return as
}

type Appointment struct {
	BaseModel
	CompanyID    uint64 `gorm:"column:CompanyID;not null"`
	ClassifyID   uint64 `gorm:"column:ClassifyID"`
	Name         string `gorm:"column:Name;not null"`
	Introduction string `gorm:"column:Introduction;default:''"` //介绍
	UseTime      string `gorm:"column:UseTime;default:''"`      //预定时间 json  {Show:true,Week:false,Begin:10,End:19,Stock:255}={是否显示,周末不显示，10点到19点}
	Orig         uint64 `gorm:"column:Orig;default:'0'"`        //原价  分
	Price        uint64 `gorm:"column:Price;default:'0'"`       //现价  分
	Link         string `gorm:"column:Link;default:''"`         //链接 json  {Show:true,Name:"dfdsf",Url:"dfdsfsdfsd"}
	//IsPayment    bool   `gorm:"column:IsPayment;default:'0'"`   //是否线上支付
	//IsPost       bool   `gorm:"column:IsPost;default:'0'"`      //是否邮寄，如果是，必须在线支付（IsPayment=true）
	Stock    uint64 `gorm:"column:DayStock;default:'0'"`  //库存
	Picture  string `gorm:"column:Picture;type:LONGTEXT"` //图片  url,url,url,url,url
	Gallery  string `gorm:"column:Gallery;type:LONGTEXT"` //显示在头部的图片 url,url,url,url,url
	Property string `gorm:"column:Property;default:'[]'"` //自定义属性 json:[{Key:"产地",Value:"福建"}]
	Params   string `gorm:"column:Params;default:'[]'"`   //购买参数 [{Key:'颜色',Value:[{P:2500,S:255,N:'黑色'},{P:2500,S:255,N:'白色'},{P:2500,S:255,N:'红色'}]}]
	//Prize        string    `gorm:"column:Prize;default:''"`        //抽奖低消费 json  {Begin:10,End:20}  10元到20元  单位 分
	//Invite    uint64    `gorm:"column:Invite;default:'0'"` //邀请好友参加，得到多少钱的优惠，分
}
type KTVRoom struct {
	BaseModel
	CompanyID      uint64  `gorm:"column:CompanyID;not null"`
	RoomName       string  `gorm:"column:RoomName;not null"`
	Characteristic string  `gorm:"column:Characteristic;default:'0'"` //特色
	Discount       float64 `gorm:"column:Discount;;default:'0'"`      //房间折扣50%
}
type KTVOrder struct {
	BaseModel
	CompanyID    uint64 `gorm:"column:CompanyID;not null"`
	ClassifyID   uint64 `gorm:"column:ClassifyID"`
	Name         string `gorm:"column:Name;not null"`
	Introduction string `gorm:"column:Introduction;default:''"` //介绍
	UseTime      string `gorm:"column:UseTime;default:''"`      //预定时间 json  {Show:true,Week:false,Begin:10,End:19,Stock:255}={是否显示,周末不显示，10点到19点}
	Orig         uint64 `gorm:"column:Orig;default:'0'"`        //原价  分
	Price        uint64 `gorm:"column:Price;default:'0'"`       //现价  分
	Link         string `gorm:"column:Link;default:''"`         //链接 json  {Show:true,Name:"dfdsf",Url:"dfdsfsdfsd"}
	//IsPayment    bool   `gorm:"column:IsPayment;default:'0'"`   //是否线上支付
	//IsPost       bool   `gorm:"column:IsPost;default:'0'"`      //是否邮寄，如果是，必须在线支付（IsPayment=true）
	Stock    uint64 `gorm:"column:DayStock;default:'0'"`  //库存
	Picture  string `gorm:"column:Picture;type:LONGTEXT"` //图片  url,url,url,url,url
	Gallery  string `gorm:"column:Gallery;type:LONGTEXT"` //显示在头部的图片 url,url,url,url,url
	Property string `gorm:"column:Property;default:'[]'"` //自定义属性 json:[{Key:"产地",Value:"福建"}]
	Params   string `gorm:"column:Params;default:'[]'"`   //购买参数 [{Key:'颜色',Value:[{P:2500,S:255,N:'黑色'},{P:2500,S:255,N:'白色'},{P:2500,S:255,N:'红色'}]}]
	//Prize        string    `gorm:"column:Prize;default:''"`        //抽奖低消费 json  {Begin:10,End:20}  10元到20元  单位 分
	//Invite    uint64    `gorm:"column:Invite;default:'0'"` //邀请好友参加，得到多少钱的优惠，分
}

func (Appointment) TableName() string {
	return "Appointment"
}

type Admin struct {
	BaseModel
	Name         string    `gorm:"column:Name"`
	Password     string    `gorm:"column:Password;not null"`
	CashPassword string    `gorm:"column:CashPassword;not null"`
	Email        string    `gorm:"column:Email;not null;unique"`
	Tel          string    `gorm:"column:Tel;not null;unique"`
	CompanyID    uint64    `gorm:"column:CompanyID;not null;unique"`
	OpenID       string    `gorm:"column:OpenID;"`
	Cash         uint64    `gorm:"column:Cash;default:'0'"`
	LastLoginAt  time.Time `gorm:"column:LastLoginAt"`
	ParentID     uint64    `gorm:"column:ParentID;not null"`
}

func (Admin) TableName() string {
	return "Admin"
}

type WxConfig struct {
	BaseModel
	CompanyID       uint64 `gorm:"column:CompanyID;unique"`
	AppID           string `gorm:"column:AppID"`
	AppSecret       string `gorm:"column:AppSecret"`
	Token           string `gorm:"column:Token"`
	EncodingAESKey  string `gorm:"column:EncodingAESKey"`
	MchID           string `gorm:"column:MchID"`
	WXOpenAppID     string `gorm:"column:WXOpenAppID"`     //[{"AppID":"wx0b83a44aa044d608","AppSecret":"ff349971f40940473bfd424aeee81b08"}]
	WXOpenAppSecret string `gorm:"column:WXOpenAppSecret"` //[{"AppID":"wx0b83a44aa044d608","AppSecret":"ff349971f40940473bfd424aeee81b08"}]
}

//wx0b83a44aa044d608 ff349971f40940473bfd424aeee81b08

func (WxConfig) TableName() string {
	return "WxConfig"
}

type Project struct {
	BaseModel
	Name     string `gorm:"column:Name;default:''"`
	Industry string `gorm:"column:Industry;default:''"`
	Budget   string `gorm:"column:Budget;default:''"`   //特色
	Describe string `gorm:"column:Describe;default:''"` //房间折扣50%
	UserName string `gorm:"column:UserName;default:''"`
	Wx       string `gorm:"column:Wx;default:''"`
	Phone    string `gorm:"column:Phone;default:''"`
	Email    string `gorm:"column:Email;default:''"`
}

func (Project) TableName() string {
	return "Project"
}

type Configuration struct {
	BaseModel
	K uint64 `gorm:"column:K;unique"`
	V string `gorm:"column:V;type:json"`
}

func (Configuration) TableName() string {
	return "Configuration"
}

type User struct {
	BaseModel
	Name        string    `gorm:"column:Name"`
	Password    string    `gorm:"column:Password;not null"`
	Tel         string    `gorm:"column:Tel;not null;unique"`
	Cash        uint64    `gorm:"column:Cash;default:'0'"`  //金额 分
	Score       uint64    `gorm:"column:Score;default:'0'"` //购物积分，100元，100个积分
	OpenID      string    `gorm:"column:OpenID"`            //
	Portrait    string    `gorm:"column:Portrait"`          //头像
	LastLoginAt time.Time `gorm:"column:LastLoginAt"`
}

func (User) TableName() string {
	return "User"
}

type Company struct {
	BaseModel
	Name         string    `gorm:"column:Name;not null"`   //店名
	Province     string    `gorm:"column:Province"`        //省
	City         string    `gorm:"column:City"`            //市
	District     string    `gorm:"column:District"`        //区域
	Address      string    `gorm:"column:Address"`         //街道地址
	Telephone    string    `gorm:"column:Telephone"`       //手机
	Categories   string    `gorm:"column:Categories"`      //门店的类型
	Longitude    string    `gorm:"column:Longitude"`       //地理位置
	Latitude     string    `gorm:"column:Latitude"`        //地理位置
	Photos       string    `gorm:"column:Photos"`          //店的图片
	Special      string    `gorm:"column:Special"`         //特色
	Opentime     string    `gorm:"column:Opentime"`        //营业时间
	Avgprice     int       `gorm:"column:Avgprice"`        //每人平均消费
	Introduction string    `gorm:"column:Introduction"`    //介绍
	Recommend    string    `gorm:"column:Recommend"`       //推荐
	Vip          int       `gorm:"column:Vip;default:'0'"` //VIP等级
	Expire       time.Time `gorm:"column:Expire"`          //过期时间
	Link         string    `gorm:"column:Link"`            //链接
}

func (Company) TableName() string {
	return "Company"
}

type Province struct {
	ID   uint64 `gorm:"column:ID;primary_key;unique"`
	P    int    `gorm:"column:P;unique"`
	Name string `gorm:"column:Name"`
}

func (Province) TableName() string {
	return "Province"
}

type City struct {
	ID   uint64 `gorm:"column:ID;primary_key;unique"`
	P    int    `gorm:"column:P"`
	C    int    `gorm:"column:C"`
	Name string `gorm:"column:Name"`
}

func (City) TableName() string {
	return "City"
}

type Area struct {
	ID   uint64 `gorm:"column:ID;primary_key;unique"`
	P    int    `gorm:"column:P"`
	C    int    `gorm:"column:C"`
	A    int    `gorm:"column:A"`
	Name string `gorm:"column:Name"`
}

func (Area) TableName() string {
	return "Area"
}

type Article struct {
	BaseModel
	Title      string `gorm:"column:Title"`
	Content    string `gorm:"column:Content;type:longtext"`
	Thumbnail  string `gorm:"column:Thumbnail"`
	CategoryID uint64 `gorm:"column:CategoryID"`
	CompanyID  uint64 `gorm:"column:CompanyID"`
	FromUrl    string `gorm:"column:FromUrl"`
}

func (Article) TableName() string {
	return "Article"
}

type Category struct {
	BaseModel
	Label string `gorm:"column:Label;unique"`
}

func (Category) TableName() string {
	return "Category"
}

type Classify struct {
	BaseModel
	Label     string `gorm:"column:Label;not null"`
	CompanyID uint64 `gorm:"column:CompanyID"`
}

func (Classify) TableName() string {
	return "Classify"
}

type ShopCoupon struct {
	BaseModel
	CompanyID uint64 `gorm:"column:CompanyID;unique"`
	Prize     string `gorm:"column:Prize;default:'{}'"` //抽奖  {Show:true,List:[{R:2000,B:10,E:20,D:7},{R:2000,B:10,E:20,D:7}]}  Require:(满多少才可以使用,0为不设条件)  10元到20元  单位 分
	//Invite uint64 `gorm:"column:Invite;default:'0'"` //邀请好友参加，得到多少钱/个的优惠，分
	//Point uint64 //返点
}

func (ShopCoupon) TableName() string {
	return "ShopCoupon"
}

type UserCoupon struct {
	BaseModel
	CompanyID uint64 `gorm:"column:CompanyID;default:'0'"`
	UserID    uint64 `gorm:"column:UserID;default:'0'"`
	Prize     string `gorm:"column:Prize;default:'[]'"` //[{R:2000,M:500,T:'2017-08-26 23:25:30',D:7},{R:2000,M:500,T:'2017-08-26 23:25:30',D:7}]
}

func (UserCoupon) TableName() string {
	return "UserCoupon"
}
