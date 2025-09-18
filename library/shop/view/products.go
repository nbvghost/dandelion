package view

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
	"strings"
)

type ProductsRequest struct {
	Organization *model.Organization `mapping:""`

	GoodsTypeUri      string   `uri:"TypeID"`
	GoodsTypeChildUri string   `uri:"SubTypeID"`
	SortName          string   `form:"sort_name"`
	Sort              string   `form:"sort"`
	PageIndex         int      `form:"page"`
	Option            []string `form:"option"`
}
type ProductsReply struct {
	extends.ViewBase
	SiteData         serviceargument.SiteData[*extends.GoodsDetail]
	OptionList       *serviceargument.Options
	SelectOptionList *serviceargument.Options
}

func (m *ProductsRequest) Render(context constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &ProductsReply{
		ViewBase: extends.ViewBase{
			Name: "products",
		},
	}

	var attributes []serviceargument.Option
	//过滤重复的参数
	for i := range m.Option {
		keys := strings.SplitN(m.Option[i], "-", 2)
		if len(keys) == 2 {
			var optionType = serviceargument.NewOptionsType(keys[0])
			var keyvalue = strings.Split(keys[1], "-")
			var key = keyvalue[0]
			var value = keyvalue[1]

			var hasKey = false
			var hasValue = false
			for ii := range attributes {
				if strings.EqualFold(string(attributes[ii].Type), string(optionType)) && strings.EqualFold(attributes[ii].Key, key) {
					hasKey = true
					subList := attributes[ii].Value
					for iii := range subList {
						if strings.EqualFold(subList[iii].Value, value) {
							hasValue = true
							break
						}
					}
				}
			}

			if hasKey == false {
				attributes = append(attributes, serviceargument.Option{
					Type:  optionType,
					Key:   key,
					Label: key, //todo
					Value: []serviceargument.OptionValue{},
				})
			}

			if hasValue == false {
				for ii := range attributes {
					if strings.EqualFold(string(attributes[ii].Type), string(optionType)) && strings.EqualFold(attributes[ii].Key, key) {
						attributes[ii].Value = []serviceargument.OptionValue{{
							Value: value,
							Count: 0,
						}}
						break
					}
				}
			}
		}
	}
	reply.SelectOptionList = &serviceargument.Options{Attributes: attributes}

	reply.SiteData = service.Site.GoodsList(context, m.Organization.ID, m.GoodsTypeUri, m.GoodsTypeChildUri, attributes, &serviceargument.SortMethod{Field: m.SortName, Method: m.Sort}, m.PageIndex, 21)

	/*optionList, err := service.Goods.ProductOptions(context, m.Organization.ID)
	if err != nil {
		return nil, err
	}*/

	newOptions := &serviceargument.Options{Attributes: make([]serviceargument.Option, 0)}
	//过滤已经选中的参数
	for i := range reply.SiteData.Options.Attributes {
		values := make([]serviceargument.OptionValue, 0)
		for ii := range reply.SiteData.Options.Attributes[i].Value {
			var has = false
			for ai := range attributes {
				if strings.EqualFold(string(attributes[ai].Type), string(reply.SiteData.Options.Attributes[i].Type)) && strings.EqualFold(attributes[ai].Key, reply.SiteData.Options.Attributes[i].Key) {
					attributes[ai].Label = reply.SiteData.Options.Attributes[i].Label
					has = true
				}
			}
			if has == false {
				values = append(values, reply.SiteData.Options.Attributes[i].Value[ii])
			}
		}
		if len(values) > 0 && len(reply.SiteData.Options.Attributes[i].Label) > 0 {
			newOptions.Attributes = append(newOptions.Attributes, serviceargument.Option{
				Type:  reply.SiteData.Options.Attributes[i].Type,
				Key:   reply.SiteData.Options.Attributes[i].Key,
				Label: reply.SiteData.Options.Attributes[i].Label,
				Value: values,
			})
		}
	}

	reply.OptionList = newOptions

	tags, err := service.Goods.Tag.FindGoodsTags(m.Organization.ID)
	if err != nil {
		return nil, err
	}
	reply.SiteData.Tags = tags

	//goodsType, goodsTypeChild := m.GoodsTypeService.GetGoodsTypeByUri(m.Organization.ID, m.GoodsTypeUri, m.GoodsTypeChildUri)

	//reply.MenusData = module.NewProductMenusData(goodsType, goodsTypeChild)

	/*menusData := m.ContentService.FindShowMenus(m.Organization.ID)
	for _, v := range menusData.List {
		if v.Type == model.ContentTypeProducts {
			reply.MenusData.Menus = v
			break
		}
	}*/

	/*pageIndex, pageSize, total, list, err := m.GoodsService.PaginationGoods(m.Organization.ID, reply.MenusData.TypeID, reply.MenusData.SubTypeID, m.PageIndex)
	if err != nil {
		return nil, err
	}
	reply.Pagination = module.NewContentPagination(pageIndex, pageSize, total, list)*/

	/*var description = reply.MenusData.Menus.Name
	for _, v := range reply.MenusData.Menus.List {
		description = description + "," + v.Name
	}

	reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := m.ContentService.GetTitle(db.Orm(), m.Organization.ID)
		meta.SetBase(reply.MenusData.Menus.Name, siteName, description)
		return nil
	}*/
	reply.HtmlMetaCallback = func(viewBase extends.ViewBase, meta *extends.HtmlMeta) error {
		siteName := service.Content.GetTitle(db.Orm(), m.Organization.ID)
		meta.SetBase(reply.SiteData.CurrentMenuData.Menus.Name, siteName, "", "")
		return nil
	}
	return reply, nil
}
