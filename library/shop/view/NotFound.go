package view

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
	"net/http"
)

type NotFound struct {
	Organization  *model.Organization  `mapping:""`
	ContentConfig *model.ContentConfig `mapping:""`
	ErrorText     string
	Stack         string
}

type NotFoundReply struct {
	extends.ViewBase
	//MenusData     module.MenusData
	//ContentConfig *model.ContentConfig
	//Organization  *model.Organization
	SiteData  serviceargument.SiteData[*model.Content]
	ErrorText string
	Stack     string
}

func (m *NotFound) Render(context constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &NotFoundReply{
		ViewBase:  extends.ViewBase{},
		ErrorText: m.ErrorText,
		Stack:     m.Stack,
	}
	contextValue := contexext.FromContext(context)
	contextValue.Response.WriteHeader(http.StatusNotFound)
	reply.Name = "404"
	reply.SiteData = service.Site.GetContentTypeByUri(context, m.Organization.ID, "", "", 0)
	return reply, nil
}
