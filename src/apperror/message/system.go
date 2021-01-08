package message

type System int

const (
	SystemError System = iota + 1
	StatusBadRequest
)

func (s System) String() string {
	var messages = map[System]string{
		SystemError:      "システムエラー",
		StatusBadRequest: "リクエストの形式が間違っています",
	}
	return messages[s]
}
