package view

import (
	"fmt"
	"github.com/nbvghost/dandelion/config"
	"io/ioutil"
	"net/http"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/result"
)

type FaviconicoRequest struct {
	ContentConfig *model.ContentConfig `mapping:""`
}

type FaviconicoReply struct {
	extends.ViewBase
	ImgBytes []byte
}

func (m *FaviconicoReply) GetResult(context constrain.IContext, viewHandler constrain.IViewHandler) constrain.IResult {
	return &result.ImageBytesResult{
		Data:        m.ImgBytes,
		ContentType: "image/x-icon",
	}
}

func (m *FaviconicoRequest) Render(context constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &FaviconicoReply{}
	reply.Name = "favicon.ico"

	contextValue := contexext.FromContext(context)

	if len(m.ContentConfig.FaviconIco) > 0 {
		ossUrl, err := context.Etcd().SelectInsideServer(config.MicroServerOSS)
		if err != nil {
			return nil, err
		}
		resp, err := http.Get(fmt.Sprintf("http://%s/assets%s", ossUrl, m.ContentConfig.FaviconIco))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		reply.ImgBytes = b
	} else {
		b, err := ioutil.ReadFile(fmt.Sprintf("view/%s/favicon.ico", contextValue.DomainName))
		if err != nil {
			return nil, err
		}
		reply.ImgBytes = b
	}
	return reply, nil
}
