package redis

import "fmt"

func NewTokenKey(token string) string { return fmt.Sprintf("token:%s", token) }
func NewUIDKey() string               { return fmt.Sprintf("uid") }
