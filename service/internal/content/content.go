package content

import (
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service/internal/company"
	"github.com/nbvghost/dandelion/service/internal/journal"
	"github.com/nbvghost/dandelion/service/internal/pinyin"
	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
)

type ContentService struct {
	model.BaseDao
	Journal             journal.JournalService
	OrganizationService company.OrganizationService
	PinyinService       pinyin.Service
}

func (service ContentService) GetTitle(orm *gorm.DB, OID dao.PrimaryKey) string {
	organization := service.OrganizationService.GetOrganization(OID).(*model.Organization)
	contentConfig := repository.ContentConfigDao.GetContentConfig(db.Orm(), organization.Primary())
	title := contentConfig.Name
	if len(title) == 0 {
		title = organization.Name
	}
	return title
}
