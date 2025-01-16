package httpext

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/server/redis"
	"time"
)

func WriteSession(context constrain.IContext, sessionID dao.PrimaryKey) error {
	expiration := time.Duration(30) * time.Minute
	err := context.Redis().Set(context, redis.NewTokenKey(context.Token()), &Session{
		ID:    fmt.Sprintf("%d", sessionID),
		Token: context.Token(),
	}, expiration)
	if err != nil {
		return err
	}
	return nil
}
