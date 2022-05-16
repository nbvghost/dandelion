package cache

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/singleton"
)

var Cache = &cacheService{
	pinyin:        map[string]string{},
	langBaiduCode: map[string]string{},
}

type cacheService struct {
	pinyin        map[string]string
	langBaiduCode map[string]string
	showLanguage  []model.Language
}

func (m *cacheService) GetPinyin(word string) string {
	return m.pinyin[word]
}

func (m *cacheService) GetLangBaiduCode() map[string]string {
	return m.langBaiduCode
}
func (m *cacheService) ShowLang() []model.Language {
	return m.showLanguage
}
func Init() {
	var cacheList []model.Pinyin
	singleton.Orm().Model(model.Pinyin{}).Find(&cacheList)
	for _, v := range cacheList {
		Cache.pinyin[v.Word] = v.Pinyin
	}

	var languageList []model.Language
	singleton.Orm().Model(model.Language{}).Where(`"CodeBiadu"<>''`).Find(&languageList)
	for index, v := range languageList {
		Cache.showLanguage = append(Cache.showLanguage, languageList[index])
		Cache.langBaiduCode[v.Code6391] = v.CodeBiadu
	}
}
