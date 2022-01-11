package action

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
