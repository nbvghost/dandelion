package notify

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"time"
)

type Login struct {
	Post struct {
		User struct {
			ID          uint64
			Account     string
			Authority   string
			LastLoginAt time.Time
		}
		AppName string
		AppKey  string
		Token   string
		Secret  string
	} `method:"Post"`
}

func (m *Login) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}

func (m *Login) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	/*SecretKey := "dfgdfgsdfgdsfggdsfgdsfgdsfgsdfgdsfgsdf"
	AppKey := "dfgfdgsghjfgdfh"

	if encryption.Md5ByString(fmt.Sprintf("%v&%v&%v", AppKey, m.Post.Token, SecretKey)) != m.Post.Secret {
		return nil, fmt.Errorf("secret验证不通过")
	}

	var admin *model.Admin
	if admin, err = service.Admin.InitOrganizationInfo(m.Post.User.Account); err != nil {
		return nil, err
	}

	b, err := json.Marshal(httpext.Session{
		ID:    fmt.Sprintf("%d", admin.ID),
		Token: m.Post.Token,
	})
	if err != nil {
		return nil, err
	}
	err = context.Redis().Set(context, redis.NewTokenKey(m.Post.Token), string(b), time.Minute*10)
	if err != nil {
		return nil, err
	}*/
	return &result.TextResult{Data: "SUCCESS"}, nil
}
