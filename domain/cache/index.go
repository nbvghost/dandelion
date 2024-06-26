package cache

import "github.com/nbvghost/dandelion/entity/model"

var Cache = struct {
	ChinesePinyinCache *ChinesePinyinCache
	LanguageCache      LanguageCache
	LanguageCodeCache  LanguageCodeCache
	RedisCache         RedisCache
}{
	ChinesePinyinCache: &ChinesePinyinCache{Pinyin: make(map[string]string)},
	LanguageCache:      LanguageCache{ShowLanguage: make([]model.Language, 0)},
	LanguageCodeCache:  LanguageCodeCache{LangBaiduCode: make(map[string]string)},
}
