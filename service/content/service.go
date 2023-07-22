package content

import (
	"github.com/nbvghost/dandelion/library/db"
	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/service/company"
	"github.com/nbvghost/dandelion/service/journal"
	"github.com/nbvghost/dandelion/service/pinyin"
)

type ContentService struct {
	model.BaseDao
	Journal             journal.JournalService
	OrganizationService company.OrganizationService
	PinyinService       pinyin.Service
}

func (service ContentService) GetTitle(orm *gorm.DB, OID dao.PrimaryKey) string {
	organization := service.OrganizationService.GetOrganization(OID).(*model.Organization)
	contentConfig := service.GetContentConfig(db.Orm(), organization.Primary())
	title := contentConfig.Name
	if len(title) == 0 {
		title = organization.Name
	}
	return title
}
