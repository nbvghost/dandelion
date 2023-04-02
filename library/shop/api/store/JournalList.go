package store

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/journal"
)

type JournalList struct {
	JournalService journal.JournalService
	Store          *model.Store `mapping:""`

	Post struct {
		StartDate string `form:"StartDate"`
		EndDate   string `form:"EndDate"`
	} `method:"Post"`
}

func (m *JournalList) HandlePost(context constrain.IContext) (constrain.IResult, error) {

	//StoreID, _ := strconv.ParseUint(context.Request.FormValue("StoreID"), 10, 64)

	list := m.JournalService.StoreListJournal(m.Store.ID, m.Post.StartDate, m.Post.EndDate)

	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: list}}, nil
}

func (m *JournalList) Handle(context constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")

}
