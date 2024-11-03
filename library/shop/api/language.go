package api

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/cache"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/environments"
	"github.com/nbvghost/dandelion/library/result"
)

type Language struct {
	Organization  *model.Organization  `mapping:""`
	ContentConfig *model.ContentConfig `mapping:""`
	Get           struct {
	} `method:"Get"`
}
type LanguageItem struct {
	Name string
	Code string
}

func (m *Language) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	contextValue := contexext.FromContext(context)
	domainName := contextValue.DomainName
	if !environments.Release() {
		domainName = fmt.Sprintf("dev.%s", domainName)
	}

	var langList []LanguageItem

	if !m.ContentConfig.EnableMultiLanguage {
		return &result.JsonResult{Data: map[string]any{"LanguageList": langList, "Language": contextValue.Lang, "DomainName": domainName}}, nil
	}
	showLang := cache.Cache.LanguageCache.ShowLang()

	for _, v := range showLang {
		langList = append(langList, LanguageItem{
			Name: v.Name,
			Code: v.Code,
		})
	}
	return &result.JsonResult{Data: map[string]any{"LanguageList": langList, "Language": contextValue.Lang, "DomainName": domainName}}, nil
}
