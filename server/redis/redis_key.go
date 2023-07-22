package redis

import (
	"fmt"

	"github.com/nbvghost/dandelion/library/dao"
)

func NewTokenKey(token string) string            { return fmt.Sprintf("token:%s", token) }
func NewUIDKey() string                          { return fmt.Sprintf("uid") }
func NewConfirmOrders(UID dao.PrimaryKey) string { return fmt.Sprintf("%d:confirm_orders", UID) }
func NewUser(UID dao.PrimaryKey) string          { return fmt.Sprintf("%d:user", UID) }
func NewArticleLookCount(UID, ArticleID dao.PrimaryKey) string {
	return fmt.Sprintf("%d:article:%d:look_count", UID, ArticleID)
}
