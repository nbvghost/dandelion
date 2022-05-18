package goods

import (
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/internal/repository"
)

type AttributesService struct {
}

func (service AttributesService) AllAttributesName() ([]*extends.GoodsAttributesNameInfo, error) {

	return repository.GoodsAttributes.QueryGoodsAttributesNameInfo()
}
func (service AttributesService) AllAttributesByName(name string) ([]*extends.GoodsAttributesValueInfo, error) {

	return repository.GoodsAttributes.QueryGoodsAttributesValueInfoByName(name)
}
