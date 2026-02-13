package widget

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/repository"
)

type CustomerService struct {
	Organization *model.Organization `mapping:""`
}

func (m *CustomerService) Template() ([]byte, error) {
	return nil, nil
}

func (m *CustomerService) Render(ctx constrain.IContext) (map[string]any, error) {
	ContentConfig := repository.ContentConfigDao.GetContentConfig(db.GetDB(ctx), m.Organization.ID)
	return map[string]any{
		"ContentConfig": ContentConfig,
	}, nil
}
