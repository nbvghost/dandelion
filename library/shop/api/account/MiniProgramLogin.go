package account

import (
	"fmt"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/shop/api/account/redisKey"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/play"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/server/httpext"
	"github.com/nbvghost/dandelion/server/redis"
	"github.com/nbvghost/dandelion/service/journal"
	"github.com/nbvghost/dandelion/service/user"
	"github.com/nbvghost/dandelion/service/wechat"
)

type MiniProgramLogin struct {
	UserService           user.UserService
	WxService             wechat.WxService
	MessageNotify         wechat.MessageNotify
	JournalService        journal.JournalService
	WXQRCodeParamsService wechat.WXQRCodeParamsService
	Post                  struct {
		Code string
		//UserInfo string
		ShareKey string
	} `method:"Post"`
	Organization *model.Organization `mapping:""`
	WechatConfig *model.WechatConfig `mapping:""`
}

func (g *MiniProgramLogin) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {

	//userInfo := make(map[string]interface{})

	//util.JSONToStruct(loginInfo.UserInfo, &userInfo)

	wxa := g.WechatConfig

	err, OpenID, SessionKey := g.WxService.MiniProgramInfo(g.Post.Code, wxa.AppID, wxa.AppSecret)
	fmt.Println(err, OpenID, SessionKey)

	if err == nil {
		tx := db.Orm().Begin()
		newUser := g.UserService.AddUserByOpenID(tx, g.Organization.ID, OpenID)
		newUser.OpenID = OpenID
		//newUser.Name = userInfo["nickName"].(string)
		//newUser.Portrait = userInfo["avatarUrl"].(string)
		//gender, _ := strconv.ParseInt(strconv.FormatFloat(userInfo["gender"].(float64), 'f', 0, 64), 10, 64)

		//newUser.Gender = int(gender)
		//newUser.OID = company.ID
		newUser.LastLoginAt = time.Now()

		if newUser.SuperiorID == 0 {
			if !strings.EqualFold(g.Post.ShareKey, "") {
				SuperiorID, _ := g.WXQRCodeParamsService.DecodeShareKey(g.Post.ShareKey)

				if newUser.ID != dao.PrimaryKey(SuperiorID) {

					//如果往上6级有包含新用户的ID，则不能绑定级别关系
					if !strings.Contains(g.UserService.LeveAll6(tx, SuperiorID), strconv.Itoa(int(newUser.ID))) {
						//var superiorUser model.User
						superiorUser := dao.GetByPrimaryKey(tx, entity.User, SuperiorID).(*model.User)
						if superiorUser.ID != 0 {
							newUser.SuperiorID = dao.PrimaryKey(SuperiorID)

							//todo
							InviteUser := 50
							err := g.JournalService.AddScoreJournal(tx,
								superiorUser.ID,
								"邀请新朋友获取积分", "邀请新朋友获取积分",
								play.ScoreJournal_Type_InviteUser, int64(InviteUser), extends.KV{Key: "SuperiorID", Value: SuperiorID})
							if err != nil {
								log.Println(err)
							}

							err = g.JournalService.AddUserJournal(tx,
								superiorUser.ID,
								"邀请新朋友获得现金", "邀请新朋友获得现金",
								play.UserJournal_Type_USER_LEVE, int64(30), extends.KV{Key: "UserID", Value: newUser.ID}, newUser.ID)
							if err != nil {
								log.Println(err)
							}

							go func(superiorUser *model.User, newUser *model.User) {
								g.MessageNotify.NewUserJoinNotify(newUser, superiorUser)
								time.Sleep(3 * time.Second)
								g.MessageNotify.INComeNotify(superiorUser, "邀请新朋友获得现金", "0小时", "收入：0.3元")
							}(superiorUser, newUser)

						}
					}

				}
			}
		}

		err = dao.UpdateByPrimaryKey(tx, entity.User, newUser.ID, newUser)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()

		err = ctx.Redis().Set(ctx, redis.NewTokenKey(ctx.Token()), &httpext.Session{
			ID:    fmt.Sprintf("%d", newUser.ID),
			Token: ctx.Token(),
		}, time.Minute*10)
		if err != nil {
			log.Println(err)
		}

		err = ctx.Redis().Set(ctx, redisKey.NewMiniProgramKey(newUser.ID), SessionKey, time.Minute*10)
		if err != nil {
			log.Println(err)
		}
		//context.Session.Attributes.Put(play.SessionUser, newUser)
		//context.Session.Attributes.Put(play.SessionOpenID, OpenID)
		//context.Session.Attributes.Put(play.SessionMiniProgramKey, SessionKey)
		//context.Session.Attributes.Put(play.SessionOrganization, company)

		//tool.CipherDecrypterData()

		results := make(map[string]interface{})
		results["User"] = newUser
		results["MyShareKey"] = g.WXQRCodeParamsService.EncodeShareKey(newUser.ID, 0) //tool.Hashids{}.Encode(newUser.ID)

		return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "登陆成功", Data: results}}, nil
	} else {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, nil
	}
}

func (g *MiniProgramLogin) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//TODO implement me
	panic("implement me")

}
