package search

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/library/dao"
)

func (Service) PutIndex(tx *gorm.DB, ID dao.PrimaryKey, Title, Content string) error {
	if err := tx.Exec(fmt.Sprintf(`UPDATE "FullTextSearch" SET "Index" = to_tsvector('english', coalesce("Title",'') || coalesce("Content",'')) WHERE "ID" = '%d'`, ID)).Error; err != nil {
		return err
	}
	return nil
}
