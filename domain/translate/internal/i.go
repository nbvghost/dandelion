package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/domain/translate/internal/aliyun"
	"github.com/nbvghost/dandelion/domain/translate/internal/volcengine"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/service"
	"log"
	"time"
)

type Translate interface {
	Translate(query []string, from, to string) (map[int]string, error)
}

type NewTranslate struct {
}

var fakeTranslate = &FakeTranslate{}

func (m *NewTranslate) checkTranslate(c *model.Configuration) [2]int {
	vs := make([]int, 0)

	tx := db.Orm()
	now := time.Now()

	_, ok := translatorConfig[string(c.K)]
	if !ok {
		return [2]int{0, 0}
	}

	err := json.Unmarshal([]byte(c.V), &vs)
	if err != nil {
		return [2]int{0, 0}
	}

	if len(vs) != 2 {
		return [2]int{0, 0}
	}

	//useWordCount := vs[0]
	//maxFreeWordCount := vs[1]

	if c.UpdatedAt.Month() != now.Month() {
		vs[0] = vs[1]
		vText, err := json.Marshal(&vs)
		if err != nil {
			return [2]int{0, 0}
		}
		err = dao.UpdateByPrimaryKey(tx, &model.Configuration{}, c.ID, map[string]any{"V": vText})
		if err != nil {
			return [2]int{0, 0}
		}
	}
	return [2]int{vs[0], vs[1]}
}
func (m *NewTranslate) Translate(query []string, from, to string) (map[int]string, error) {
	var translateWordCount int
	for _, s := range query {
		translateWordCount = translateWordCount + len(s)
	}

	for key := range translatorConfig {

		translator := translatorConfig[key]

		config := service.Configuration.GetTranslate(db.Orm(), key)
		if config.IsZero() {
			return nil, errors.New("没有配制翻译参数") //fakeTranslate.Translate(query, from, to)
		}

		wordCounts := m.checkTranslate(config)

		ableWordCount := wordCounts[0]
		//maxFreeWordCount := wordCounts[1]

		if translateWordCount < ableWordCount {

			translateData, err := translator.Translate(query, from, to)
			if err != nil {
				log.Println(fmt.Sprintf("翻译器[%s]翻译出错，使用下一下翻译器。%s", config.K, err.Error()))
				continue
			}
			wordCounts[0] = ableWordCount - translateWordCount
			err = dao.UpdateByPrimaryKey(db.Orm(), &model.Configuration{}, config.ID, map[string]any{"V": util.StructToJSON(wordCounts)})
			if err != nil {
				return nil, err
			}
			return translateData, nil

		}
	}
	return nil, errors.New("没有可用的翻译器")
}

type FakeTranslate struct{}

func (m *FakeTranslate) Translate(query []string, from, to string) (map[int]string, error) {
	d := make(map[int]string)
	for i, s := range query {
		d[i] = s
	}
	return d, nil
}

var translatorConfig = map[string]Translate{
	//"TranslateBaidu":      baidu.New(),
	"TranslateAliyun":     aliyun.New(),
	"TranslateVolcengine": volcengine.New(),
	//{K: "Libre", MaxFreeWordCount: 2000000 - 100, T: libre.New()},
}

func New() (Translate, error) {
	return &NewTranslate{}, nil
}
