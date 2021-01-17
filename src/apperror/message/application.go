package message

type Application int

const (
	SystemError Application = iota + 1
	StatusBadRequest
	WrongPassword
	LoginFailed
)

func (s Application) String() string {
	var messages = map[Application]string{
		SystemError:      "システムエラー",
		StatusBadRequest: "リクエストの形式に誤りがあります",
		WrongPassword:    "パスワードが間違っています",
		LoginFailed:      "ログインに失敗しました",
	}
	return messages[s]
}
