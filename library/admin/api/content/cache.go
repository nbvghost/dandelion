package content

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/repository"
	"io"
	"net/http"
)

type CacheAction string

const (
	CacheActionClear CacheAction = "clear"
)

type Cache struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		Action CacheAction
	} `method:"Post"`
}

func (m *Cache) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	return nil, nil
}

func (m *Cache) HandlePost(context constrain.IContext) (constrain.IResult, error) {
	dns := repository.DNSDao.GetDefaultDNS(m.Organization.ID)
	if dns.IsZero() {
		return nil, result.NewErrorText("找不到DNS记录")
	}
	server, err := context.Etcd().SelectInsideServer(&config.MicroServer{Name: dns.Domain, ServerType: config.ServerTypeHttp})
	if err != nil {
		return nil, err
	}
	params := map[string]any{"Action": m.Post.Action, "OID": m.Organization.ID}

	paramsJSON, err := json.Marshal(&params)
	if err != nil {
		return nil, err
	}

	post, err := http.Post(fmt.Sprintf("http://%s/api/inside/cache", server), "application/json", bytes.NewReader(paramsJSON))
	if err != nil {
		return nil, err
	}
	defer post.Body.Close()
	body, err := io.ReadAll(post.Body)
	if err != nil {
		return nil, err
	}
	var action result.ActionResult
	err = json.Unmarshal(body, &action)
	if err != nil {
		return nil, err
	}
	return &action, nil
}
