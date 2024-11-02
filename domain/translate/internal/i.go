package internal

import (
	"fmt"
	"github.com/nbvghost/dandelion/domain/translate/internal/aliyun"
	"github.com/nbvghost/dandelion/domain/translate/internal/baidu"
	"github.com/nbvghost/dandelion/domain/translate/internal/volcengine"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/tool/object"
	"gorm.io/gorm"
	"log"
	"time"
)

type Translate interface {
	Translate(query []string, from, to string) (map[int]string, error)
}

type NewTranslate struct {
}

var fakeTranslate = &FakeTranslate{}

func (m *NewTranslate) Translate(query []string, from, to string) (map[int]string, error) {
	var err error
	tx := db.Orm().Begin()
	defer func() {
		if e := recover(); e != nil {
			tx.Rollback()
			return
		}
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	ts := service.Configuration.GetTranslate(tx, 0)
	if len(ts) == 0 {
		return fakeTranslate.Translate(query, from, to)
	}

	var translateWordCount int
	for _, s := range query {
		translateWordCount = translateWordCount + len(s)
	}

	now := time.Now()
	for i, t := range ts {
		tc, ok := translatorConfig[string(t.K)]
		if !ok {
			tc = &TranslatorConfig{
				K:                string(t.K),
				MaxFreeWordCount: 0,
				T:                fakeTranslate,
			}
		}
		if t.UpdatedAt.Month() != now.Month() {
			ts[i].V = fmt.Sprintf("%d", tc.MaxFreeWordCount)
			err = dao.UpdateByPrimaryKey(tx, &model.Configuration{}, ts[i].ID, map[string]any{"V": ts[i].V})
			if err != nil {
				return nil, err
			}
		}
	}

	var currentTranslator *TranslatorConfig

	var id dao.PrimaryKey
	for i := range ts {
		tc := translatorConfig[string(ts[i].K)]
		ableCount := object.ParseInt(ts[i].V)
		if translateWordCount < ableCount {
			currentTranslator = tc
			id = ts[i].ID
			break
		}
	}

	if currentTranslator == nil {
		log.Printf("无法获取翻译器，请查询配制:%+v\n", translatorConfig)
		return fakeTranslate.Translate(query, from, to)
	}

	err = dao.UpdateByPrimaryKey(tx, &model.Configuration{}, id, map[string]any{"V": gorm.Expr(fmt.Sprintf(`"V"::int - ?`), translateWordCount)})
	if err != nil {
		return nil, err
	}
	var d map[int]string
	d, err = currentTranslator.T.Translate(query, from, to)
	return d, err

	//return libre.New().Translate(query, from, to)
}

type FakeTranslate struct{}

func (m *FakeTranslate) Translate(query []string, from, to string) (map[int]string, error) {
	d := make(map[int]string)
	for i, s := range query {
		d[i] = s
	}
	return d, nil
}

type TranslatorConfig struct {
	K                string
	MaxFreeWordCount int
	T                Translate
}

var translatorConfig = map[string]*TranslatorConfig{
	"TranslateBaidu":      {K: "TranslateBaidu", MaxFreeWordCount: 2000000 - 100, T: baidu.New()},
	"TranslateAliyun":     {K: "TranslateAliyun", MaxFreeWordCount: 1000000 - 100, T: aliyun.New()},
	"TranslateVolcengine": {K: "TranslateVolcengine", MaxFreeWordCount: 2000000 - 100, T: volcengine.New()},
	//"Libre": {K: "Libre", MaxFreeWordCount: 2000000 - 100, T: libre.New()},
}

func New() (Translate, error) {
	return &NewTranslate{}, nil
}
