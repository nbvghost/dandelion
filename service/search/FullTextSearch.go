package search

import (
	"fmt"
	"github.com/nbvghost/dandelion/library/db"

	"github.com/nbvghost/gpa/types"
)

func (Service) PutIndex(ID types.PrimaryKey, Title, Content string) error {
	if err := db.Orm().Exec(fmt.Sprintf(`UPDATE "FullTextSearch" SET "Index" = to_tsvector('english', coalesce("Title",'') || coalesce("Content",'')) WHERE "ID" = '%d'`, ID)).Error; err != nil {
		return err
	}
	return nil
}
