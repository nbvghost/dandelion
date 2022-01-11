package redis

import "fmt"

func NewTokenKey(token string) string { return fmt.Sprintf("token:%s", token) }
