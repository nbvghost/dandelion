package entity

import (
	"fmt"
	"reflect"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
)

var (
	User                 = Register((*model.User)(nil))
	UserInfo             = Register((*model.UserInfo)(nil))
	Admin                = Register((*model.Admin)(nil))
	Configuration        = Register((*model.Configuration)(nil))
	Logger               = Register((*model.Logger)(nil))
	UserJournal          = Register((*model.UserJournal)(nil))
	StoreJournal         = Register((*model.StoreJournal)(nil))
	CollageRecord        = Register((*model.CollageRecord)(nil))
	CollageGoods         = Register((*model.CollageGoods)(nil))
	District             = Register((*model.District)(nil))
	SupplyOrders         = Register((*model.SupplyOrders)(nil))
	StoreStock           = Register((*model.StoreStock)(nil))
	Verification         = Register((*model.Verification)(nil))
	FullCut              = Register((*model.FullCut)(nil))
	Voucher              = Register((*model.Voucher)(nil))
	TimeSell             = Register((*model.TimeSell)(nil))
	Store                = Register((*model.Store)(nil))
	ExpressTemplate      = Register((*model.ExpressTemplate)(nil))
	CardItem             = Register((*model.CardItem)(nil))
	Goods                = Register((*model.Goods)(nil))
	GoodsType            = Register((*model.GoodsType)(nil))
	GoodsTypeChild       = Register((*model.GoodsTypeChild)(nil))
	OrdersGoods          = Register((*model.OrdersGoods)(nil))
	ScoreJournal         = Register((*model.ScoreJournal)(nil))
	ScoreGoods           = Register((*model.ScoreGoods)(nil))
	Specification        = Register((*model.Specification)(nil))
	Transfers            = Register((*model.Transfers)(nil))
	Orders               = Register((*model.Orders)(nil))
	ShoppingCart         = Register((*model.ShoppingCart)(nil))
	Rank                 = Register((*model.Rank)(nil))
	GiveVoucher          = Register((*model.GiveVoucher)(nil))
	Organization         = Register((*model.Organization)(nil))
	OrganizationJournal  = Register((*model.OrganizationJournal)(nil))
	OrdersPackage        = Register((*model.OrdersPackage)(nil))
	Manager              = Register((*model.Manager)(nil))
	Collage              = Register((*model.Collage)(nil))
	Content              = Register((*model.Content)(nil))
	ContentItem          = Register((*model.ContentItem)(nil))
	ContentType          = Register((*model.ContentType)(nil))
	ContentSubType       = Register((*model.ContentSubType)(nil))
	Question             = Register((*model.Question)(nil))
	QuestionTag          = Register((*model.QuestionTag)(nil))
	AnswerQuestion       = Register((*model.AnswerQuestion)(nil))
	TimeSellGoods        = Register((*model.TimeSellGoods)(nil))
	ContentConfig        = Register((*model.ContentConfig)(nil))
	WXQRCodeParams       = Register((*model.WXQRCodeParams)(nil))
	GoodsAttributes      = Register((*model.GoodsAttributes)(nil))
	GoodsAttributesGroup = Register((*model.GoodsAttributesGroup)(nil))
	GoodsWish            = Register((*model.GoodsWish)(nil))
	LeaveMessage         = Register((*model.LeaveMessage)(nil))
	Subscribe            = Register((*model.Subscribe)(nil))
	FullTextSearch       = Register((*model.FullTextSearch)(nil))
	Pinyin               = Register((*model.Pinyin)(nil))
	Language             = Register((*model.Language)(nil))
	Translate            = Register((*model.Translate)(nil))
	DNS                  = Register((*model.DNS)(nil))
	Advert               = Register((*model.Advert)(nil))
	WechatConfig         = Register((*model.WechatConfig)(nil))
	PushData             = Register((*model.PushData)(nil))
	ExpressCompany       = Register((*model.ExpressCompany)(nil))
	Area                 = Register((*model.Area)(nil))
)

func GetModel(name string) (dao.IEntity, error) {
	return defaultModel.GetModel(name)
}
func Register(e dao.IEntity) dao.IEntity {
	return defaultModel.Register(e)
}

var defaultModel = New()

type r struct {
	models map[string]dao.IEntity
}

func (m *r) Register(e dao.IEntity) dao.IEntity {
	name := reflect.TypeOf(e).Elem().Name()
	if _, ok := m.models[name]; ok {
		panic(fmt.Errorf("model %s 已经注册", e.TableName()))
	} else {
		m.models[name] = e
		return e
	}
}
func (m *r) GetModel(name string) (dao.IEntity, error) {
	if k, ok := m.models[name]; ok {
		return reflect.New(reflect.ValueOf(k).Elem().Type()).Interface().(dao.IEntity), nil
	} else {
		return nil, fmt.Errorf("model %s 不存在", name)
	}
}
func New() *r {
	return &r{models: map[string]dao.IEntity{}}
}
