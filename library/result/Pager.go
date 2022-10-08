package result

import "math"

type Pager struct {
	Data        interface{}
	Total       int
	Limit       int
	Offset      int //0---
	OffsetCount int
}

func (p *Pager) Calculation() Pager {

	OffsetCount := int(math.Ceil(float64(p.Total) / float64(p.Limit)))

	return Pager{
		Data:        p.Data,
		Total:       p.Total,
		Limit:       p.Limit,
		Offset:      p.Offset,
		OffsetCount: OffsetCount,
	}

}

type Pagination struct {
	List     any //这个字段要改成data 表示，可以用对象或数组
	Total    int
	PageSize int
	Page     int
}

func NewPagination(page, pageSize, total int, list any) *Pagination {
	return &Pagination{
		List:     list,
		Total:    total,
		PageSize: pageSize,
		Page:     page,
	}
}
