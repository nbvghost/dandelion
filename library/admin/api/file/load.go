package file

import (
	"fmt"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"io/ioutil"
	"net/http"
)

type Load struct {
	Get struct {
		Path string `form:"path"`
	} `method:"GET"`
}

func (m *Load) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	ossHost, err := context.Etcd().SelectInsideServer(config.MicroServerOSS)
	if err != nil {
		return nil, err
	}

	imgUrl := fmt.Sprintf("http://%s/assets%s", ossHost, m.Get.Path)

	var post *http.Response
	post, err = http.Get(imgUrl)
	if err != nil {
		return nil, err
	}
	defer post.Body.Close()

	body, err := ioutil.ReadAll(post.Body)
	if err != nil {
		return nil, err
	}

	return &result.ImageBytesResult{Data: body}, nil

}
