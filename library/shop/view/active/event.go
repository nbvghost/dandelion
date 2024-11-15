package active

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
)

type Event struct {
	Organization  *model.Organization `mapping:""`
	ContentConfig model.ContentConfig `mapping:""`
	Event         string              `uri:"event"`
}
type EventReply struct {
	extends.ViewBase
}

func (m *Event) Render(context constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &EventReply{}
	return reply, nil
}
