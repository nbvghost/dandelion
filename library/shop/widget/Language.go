package widget

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/cache"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/environments"
)

type Language struct {
	Organization  *model.Organization  `mapping:""`
	ContentConfig *model.ContentConfig `mapping:""`
}

func (m *Language) Template() ([]byte, error) {
	return nil, nil
}

type LanguageItem struct {
	Name    string
	ISOCode string
}

func (m *Language) Render(context constrain.IContext) (map[string]any, error) {
	contextValue := contexext.FromContext(context)
	domainName := contextValue.DomainName
	if !environments.Release() {
		domainName = fmt.Sprintf("dev.%s", domainName)
	}

	var langList []LanguageItem

	if !m.ContentConfig.EnableMultiLanguage {
		return map[string]any{"LanguageList": langList, "Language": contextValue.Lang, "DomainName": domainName}, nil
	} else {
		showLang := cache.Cache.LanguageCache.ShowLang()
		for _, v := range showLang {
			langList = append(langList, LanguageItem{
				Name:    v.Name,
				ISOCode: v.Code,
			})
		}
	}
	return map[string]any{"LanguageList": langList, "Language": contextValue.Lang, "DomainName": domainName}, nil
}
