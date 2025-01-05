package api

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/translate"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"strings"
)

type Translate struct {
	Post struct {
		Query []string
		From  string
		To    string
	} `method:"Post"`
}

func (m *Translate) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	return &result.JsonResult{Data: result.NewSuccess("OK")}, err
}
func (m *Translate) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	tran, err := translate.NewTranslate()
	if err != nil {
		return nil, err
	}

	var translateModelList []model.Translate
	db.Orm().Model(model.Translate{}).Where(`"LangType"=? and "Text" in ?`, m.Post.To, m.Post.Query).Find(&translateModelList)

	outList := make([]string, 0)

	for _, q := range m.Post.Query {
		var hasIndex = -1
		for i := range translateModelList {
			if strings.EqualFold(translateModelList[i].Text, q) {
				hasIndex = i
				break
			}
		}
		if hasIndex >= 0 {
			outList = append(outList, translateModelList[hasIndex].LangText)
		} else {

			translateText, err := tran.Translate([]string{q}, m.Post.From, m.Post.To)
			if err != nil {
				return nil, err
			}

			err = dao.Create(db.Orm(), &model.Translate{
				Text:     q,
				LangType: m.Post.To,
				LangText: translateText[0],
			})
			if err != nil {
				return nil, err
			}
			outList = append(outList, translateText[0])
		}
	}

	return result.NewData(map[string]any{"List": outList}), err
}
