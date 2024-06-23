package cache

import (
	"github.com/nbvghost/dandelion/entity/model"
)

type ChinesePinyinCache struct {
	Pinyin map[string]string
}

func (m *ChinesePinyinCache) GetPinyin(word string) string {
	return m.Pinyin[word]
}

type LanguageCache struct {
	ShowLanguage []model.Language
}

func (m *LanguageCache) ShowLang() []model.Language {
	return m.ShowLanguage
}

type LanguageCodeCache struct {
	LangBaiduCode map[string]string
}
func (m *LanguageCodeCache) GetLangBaiduCodeByCode6391(Code6391 string) string {
	return m.LangBaiduCode[Code6391]
}
func (m *LanguageCodeCache) GetLangBaiduCode() map[string]string {
	return m.LangBaiduCode
}
