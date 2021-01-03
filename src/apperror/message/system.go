package message

type System int

const (
	SystemError System = iota + 1
)

func (s System) String() string {
	var messages = map[System]string{
		SystemError: "システムエラー",
	}
	return messages[s]
}
