package redis

import (
	"fmt"

	"github.com/nbvghost/gpa/types"
)

func NewTokenKey(token string) string              { return fmt.Sprintf("token:%s", token) }
func NewUIDKey() string                            { return fmt.Sprintf("uid") }
func NewConfirmOrders(UID types.PrimaryKey) string { return fmt.Sprintf("%d:confirm_orders", UID) }
func NewUser(UID types.PrimaryKey) string          { return fmt.Sprintf("%d:user", UID) }
func NewArticleLookCount(UID, ArticleID types.PrimaryKey) string {
	return fmt.Sprintf("%d:article:%d:look_count", UID, ArticleID)
}
