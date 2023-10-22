package journal

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/user"
)

type ListJournal struct {
	UserService user.UserService
	User        *model.User `mapping:""`
}

func (g *ListJournal) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}

func (g *ListJournal) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)
	///CarID, _ := strconv.ParseUint(context.PathParams["CarID"], 10, 64)
	//var list []model.UserJournal
	Orm := db.Orm()
	//g.UserService.FindOrderWhere(Orm, `"CreatedAt" desc`, &list, &model.UserJournal{UserID: g.User.ID})

	list := dao.Find(Orm, &model.UserJournal{}).Where(`"UserID"=?`, g.User.ID).Order(`"CreatedAt" desc`).List()

	//err := controller.Car.FindWhere(entity.Orm(), &list, &entity.CarRecord{CarID: CarID})
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", list)}, nil

}
