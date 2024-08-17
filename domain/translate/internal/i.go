package internal

import "github.com/nbvghost/dandelion/domain/translate/internal/aliyun"

type Translate interface {
	Translate(query []string, from, to string) (map[int]string, error)
}

func New() (Translate, error) {

	return aliyun.New(), nil
}
