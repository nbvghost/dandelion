package redisKey

import (
	"fmt"
	"github.com/nbvghost/gpa/types"
)

func NewMiniProgramKey(UID types.PrimaryKey) string {
	return fmt.Sprintf("%d:mini-program-key", UID)
}
