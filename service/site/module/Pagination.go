package module

import "github.com/nbvghost/dandelion/entity/model"

type Paging struct {
	Index   int
	Disable bool
}

type ListType interface {
	~*model.Content | ~*model.FullTextSearch | ~*model.Goods | ~*model.ContentItem
}
type Pagination[T ListType] struct {
	PageIndex int
	PageSize  int
	Total     int
	List      []T
	Pages     []int

	Previous Paging
	Next     Paging
	First    Paging
	End      Paging
}

func NewPagination[T ListType](pageIndex, pageSize, total int, list []T) Pagination[T] {
	pagination := Pagination[T]{}
	pagination.PageIndex = pageIndex
	pagination.PageSize = pageSize
	pagination.Total = total
	pagination.List = list

	{
		var pages []int
		p := total / pageSize
		for i := 0; i < p; i++ {
			pages = append(pages, i)
		}
		if total%pageSize > 0 {
			pages = append(pages, p)
		}
		pagination.Pages = pages
	}

	pagination.Previous.Index = pageIndex - 1
	pagination.Next.Index = pageIndex + 1
	pagination.First.Index = 0
	if len(pagination.Pages) > 0 {
		pagination.End.Index = pagination.Pages[len(pagination.Pages)-1]
	} else {
		pagination.End.Index = 0
	}

	if pagination.PageIndex == 0 {
		pagination.Previous.Disable = true
		pagination.First.Disable = true

		pagination.Previous.Index = 0
		pagination.First.Index = 0
	}

	if pagination.PageIndex >= pagination.End.Index {
		pagination.Next.Disable = true
		pagination.End.Disable = true

		pagination.Next.Index = pagination.End.Index

	}

	return pagination
}
