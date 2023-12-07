package key

import "fmt"

func NewRedisVerifyPhoneCodeKey(phone string) string {
	return fmt.Sprintf("verify-phone-code:%s", phone)
}
func NewRedisVerifyPhoneIPKey(ip string) string {
	return fmt.Sprintf("verify-phone-ip:%s", ip)
}
