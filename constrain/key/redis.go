package key

import (
	"fmt"
	"github.com/nbvghost/dandelion/library/dao"
)

func NewRedisVerifyPhoneCodeKey(phone string) string {
	return fmt.Sprintf("verify-phone-code:%s", phone)
}
func NewRedisVerifyPhoneIPKey(ip string) string {
	return fmt.Sprintf("verify-phone-ip:%s", ip)
}
func NewMiniProgramRedisKey(UID dao.PrimaryKey) string {
	return fmt.Sprintf("%d:mini-program-key", UID)
}
func NewPaypalAccessTokenRedisKey(oid dao.PrimaryKey) string {
	return fmt.Sprintf("payment:paypal:%d:access-token", oid)
}
