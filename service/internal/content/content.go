package content

import (
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service/internal/company"
	"github.com/nbvghost/dandelion/service/internal/journal"
	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
)

type ContentService struct {
	model.BaseDao
	Journal             journal.JournalService
	OrganizationService company.OrganizationService
}

func (m ContentService) GetTitle(tx *gorm.DB, OID dao.PrimaryKey) string {
	organization := m.OrganizationService.GetOrganization(tx, OID).(*model.Organization)
	contentConfig := repository.ContentConfigDao.GetContentConfig(tx, organization.Primary())
	title := contentConfig.Name
	if len(title) == 0 {
		title = organization.Name
	}
	return title
}

func (m ContentService) GetByTitle(orm *gorm.DB, OID dao.PrimaryKey, title string) *model.Content {
	return dao.GetBy(orm, &model.Content{}, map[string]any{"UseType": "tag", "Title": title, "OID": OID}).(*model.Content)
}
