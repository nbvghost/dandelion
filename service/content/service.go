package content

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/company"
	"github.com/nbvghost/dandelion/service/journal"
	"github.com/nbvghost/dandelion/service/pinyin"
	"github.com/nbvghost/gpa/types"
	"gorm.io/gorm"
)

type ContentService struct {
	model.BaseDao
	Journal             journal.JournalService
	OrganizationService company.OrganizationService
	PinyinService       pinyin.Service
}

func (service ContentService) GetTitle(orm *gorm.DB, OID types.PrimaryKey) string {
	organization := service.OrganizationService.GetOrganization(OID)
	contentConfig := service.GetContentConfig(singleton.Orm(), organization.ID)
	title := contentConfig.Name
	if len(title) == 0 {
		title = organization.Name
	}
	return title
}
