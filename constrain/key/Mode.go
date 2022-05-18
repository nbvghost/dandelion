package key

type Mode string

const (
	ModeDev     Mode = "dev"
	ModeRelease Mode = "release"
)

func (m Mode) String() string {
	return string(m)
}
