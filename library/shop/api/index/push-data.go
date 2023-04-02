package index

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/goods"
)

type PushData struct {
	GoodsService goods.GoodsService
	Organization *model.Organization `mapping:""`
	Post         any                 `method:"post"`
}

func (m *PushData) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	return nil, nil
}
func (m *PushData) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	body, err := json.Marshal(m.Post)
	if err != nil {
		return nil, err
	}

	singleton.Orm().Model(model.PushData{}).Create(&model.PushData{
		Content: string(body),
	})
	post, err := http.Post("https://admin.7846.com/api/weiXin1Callback", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer post.Body.Close()
	body, err = ioutil.ReadAll(post.Body)
	if err != nil {
		return nil, err
	}
	return result.NewData(nil), nil
}
