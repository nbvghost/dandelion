package redisKey

import (
	"fmt"
	"github.com/nbvghost/dandelion/library/dao"
)

func NewMiniProgramKey(UID dao.PrimaryKey) string {
	return fmt.Sprintf("%d:mini-program-key", UID)
}
