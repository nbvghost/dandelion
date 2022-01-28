package environments

import "flag"

var env environments

type environments struct {
	Release bool
}

func init() {
	flag.BoolVar(&env.Release, "release", true, "release")
}
func Release() bool {
	return env.Release
}
