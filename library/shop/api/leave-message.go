package api

import (
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/service"
	"strings"
)

type LeaveMessage struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		Name          string               `form:"Name"`
		Email         string               `form:"Email"`
		SocialType    []sqltype.SocialType `form:"SocialType"`
		SocialAccount []string             `form:"SocialAccount"`
		Content       string               `form:"Content"`
	} `method:"Post"`
}

func (m *LeaveMessage) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *LeaveMessage) HandlePost(context constrain.IContext) (constrain.IResult, error) {
	contextValue := contexext.FromContext(context)
	leaveMessage := &model.LeaveMessage{}
	leaveMessage.OID = m.Organization.ID
	leaveMessage.Name = m.Post.Name
	leaveMessage.Email = m.Post.Email
	/*leaveMessage.SocialAccount = sqltype.SocialAccount{
		Type:    m.Post.SocialType,
		Hide:    false,
		Account: m.Post.SocialAccount,
	}*/
	leaveMessage.Content = m.Post.Content
	leaveMessage.ClientIP = util.GetIP(contextValue.Request)

	extend := map[string]interface{}{}
	_ = contextValue.Request.ParseForm()
	for k := range contextValue.Request.Form {
		if strings.EqualFold(k, "Name") || strings.EqualFold(k, "Email") || strings.EqualFold(k, "Content") || strings.EqualFold(k, "ClientIP") {
			continue
		}
		extend[k] = contextValue.Request.Form.Get(k)
	}

	extend["Location"], _ = util.GetIPLocation(leaveMessage.ClientIP)

	extendByte, _ := json.Marshal(&extend)
	leaveMessage.Extend = string(extendByte)

	err := db.Orm().Model(model.LeaveMessage{}).Create(leaveMessage).Error
	if err != nil {
		return nil, err
	}

	{
		botMessage := strings.Builder{}
		botMessage.WriteString(fmt.Sprintf("网站[%s]有新的消息，请注意查收。\n", context.AppName()))
		botMessage.WriteString(fmt.Sprintf(">姓名:%s\n", leaveMessage.Name))
		botMessage.WriteString(fmt.Sprintf(">Email:%s\n", leaveMessage.Email))
		botMessage.WriteString(fmt.Sprintf(">内容:%s\n", leaveMessage.Content))
		botMessage.WriteString(fmt.Sprintf(">IP:%s\n", leaveMessage.ClientIP))
		for k := range extend {
			botMessage.WriteString(fmt.Sprintf(">%s:%s\n", k, extend[k]))
		}
		err = service.Wechat.SendText(botMessage.String())
		if err != nil {
			return nil, err
		}
	}
	return &result.JsonResult{Data: result.ActionResult{}}, err
}
