package service

import (
	"dandelion/app/service/dao"
	"errors"
	"strings"
)

type ClassifyService struct {
	dao.ClassifyDao
}

func (self ClassifyService) AddClassifyNotNull(classify *dao.Classify) error {
	if strings.EqualFold(classify.Label, "") {
		return errors.New("分类名称不能空字符")
	}
	return self.Add(Orm, classify)
}
