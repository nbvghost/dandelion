package key

import (
	"golang.org/x/exp/rand"
	"time"
)

var Random = rand.New(rand.NewSource(uint64(time.Now().Unix())))
