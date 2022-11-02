package pagination

import (
	"math"

	"github.com/nbvghost/tool/object"
)

type SwitchPage struct {
	Pages    []int
	PageNum  int
	Page     int
	Total    int
	Size     int
	PrevPage int
	NextPage int
}

func GetSwitchPage(index, total int, size int) SwitchPage {
	pageNum := object.ParseInt(math.Ceil(float64(total) / float64(size)))

	prevPage := (index + 1) - 1
	if prevPage <= 0 {
		prevPage = -1
	}
	nextPage := (index + 1) + 1
	if nextPage > pageNum {
		nextPage = -1
	}

	pages := make([]int, 0)
	showPageNum := 6

	starIndex := index - (showPageNum / 2)
	endIndex := index + (showPageNum / 2)

	if starIndex <= 0 || index < (showPageNum/2) {
		for i := 0; i < showPageNum; i++ {
			if i+1 <= pageNum {
				pages = append(pages, i+1)
			}
		}
		if len(pages) == 0 {
			return SwitchPage{Pages: pages, PageNum: pageNum, Page: index + 1, Total: total, Size: size, PrevPage: prevPage, NextPage: nextPage}
		}
		if pages[0] >= 2 {
			pages = append([]int{1, -1}, pages[:]...)
		}

		if pageNum-pages[len(pages)-1] >= 2 {
			pages = append(pages, -1)
			pages = append(pages, pageNum)
		} else {
			if pages[len(pages)-1] != pageNum {
				pages = append(pages, pageNum)
			}

		}
	} else {

		for i := starIndex + 1; i < endIndex+1; i++ {
			if i <= pageNum {
				pages = append(pages, i)
			}
		}
		if len(pages) == 0 {
			return SwitchPage{Pages: pages, PageNum: pageNum, Page: index + 1, Total: total, Size: size, PrevPage: prevPage, NextPage: nextPage}
		}
		if pages[0]-1 >= 2 {
			pages = append([]int{1, -1}, pages[:]...)
		} else {
			pages = append([]int{1}, pages[:]...)
		}

		if pageNum-pages[len(pages)-1] >= 2 {
			pages = append(pages, -1)
			pages = append(pages, pageNum)
		}
	}
	return SwitchPage{Pages: pages, PageNum: pageNum, Page: index + 1, Total: total, Size: size, PrevPage: prevPage, NextPage: nextPage}
}
