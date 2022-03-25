package content

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/service/journal"
)

type ContentService struct {
	model.BaseDao
	Journal journal.JournalService
}
