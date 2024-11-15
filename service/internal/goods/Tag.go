package goods

import (
	"fmt"
	"github.com/nbvghost/dandelion/domain/tag"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
)

type TagService struct {
}

func (service TagService) FindGoodsTags(OID dao.PrimaryKey) ([]extends.Tag, error) {
	//SELECT unnest("Tags") as Tag,count("Tags") as Count FROM "Content" where  group by unnest("Tags");
	var tags []extends.Tag
	err := db.Orm().Model(model.Goods{}).Select(`unnest("Tags") as "Name",count("Tags") as "Count"`).Where(map[string]interface{}{
		"OID": OID,
	}).Where(`array_length("Tags",1)>0`).Group(`unnest("Tags")`).Order(`"Count" desc`).Find(&tags).Error
	tags = tag.CreateUri(tags)
	return tags, err
}
func (service TagService) FindGoodsByTag(OID dao.PrimaryKey, tag extends.Tag, _pageIndex int, orders ...dao.Sort) (pageIndex, pageSize int, total int64, list []*model.Goods, err error) {
	//select * from "Content" where array_length("Tags",1) is null;
	db := db.Orm().Model(model.Goods{}).Where(`"OID"=?`, OID).
		Where(`array_length("Tags",1) is not null`).
		Where(`"Tags" @> array[?]`, tag.Name)

	db.Count(&total)

	for _, v := range orders {
		db.Order(fmt.Sprintf(`"%s" %s`, v.ColumnName, v.Method))
	}

	pageSize = 20

	err = db.Limit(pageSize).Offset(_pageIndex * pageSize).Find(&list).Error
	pageIndex = _pageIndex

	return
}
