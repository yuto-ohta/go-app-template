package messages

type System int

const (
	SystemError System = iota
)

func (e System) String() string {
	switch e {
	case SystemError:
		fallthrough
	default:
		return "システムエラー"
	}
}
