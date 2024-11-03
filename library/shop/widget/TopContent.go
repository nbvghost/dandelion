package widget

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service"
)

type TopContent struct {
	Organization *model.Organization `mapping:""`

	Style        string                `arg:""`
	SoreField    string                `arg:""` //排序字段
	SoreMethod   int                   `arg:""` //排序字段
	Type         model.ContentTypeType `arg:""`
	TemplateName string                `arg:""`
	Count        uint                  `arg:""`
}

func (m *TopContent) Template() ([]byte, error) {
	return nil, nil
}

func (m *TopContent) Render(ctx constrain.IContext) (map[string]any, error) {

	contentItem := repository.ContentItemDao.GetContentItemByTypeTemplateName(db.Orm(), m.Organization.ID, m.Type, m.TemplateName)

	var list = repository.ContentDao.SortList(m.Organization.ID, contentItem.ID, m.SoreField, m.SoreMethod, m.Count)
	/*if m.TopType == TopTypeContentView {
		list = m.ContentService.HotViewList(m.Organization.ID, contentItem.ID, m.Count)
	} else if m.TopType == TopTypeContentLike {
		list = m.ContentService.HotLikeList(m.Organization.ID, contentItem.ID, m.Count)
	}*/

	var templateName string
	menusData := service.Site.FindAllMenus(m.Organization.ID)
	for i := 0; i < len(menusData.List); i++ {
		if menusData.List[i].ID == contentItem.ID {
			templateName = menusData.List[i].TemplateName
		}
	}
	return map[string]any{
		"List":         list,
		"Style":        m.Style,
		"SoreField":    m.SoreField,
		"SoreMethod":   m.SoreMethod,
		"TemplateName": templateName,
		"ContentItem":  contentItem,
	}, nil
}
