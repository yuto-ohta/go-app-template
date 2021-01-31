package message

type Application int

const (
	SystemError Application = iota + 1
	StatusBadRequest
	WrongPassword
	LoginFailed
	LogoutFailed
	UnAuthorized
	Forbidden
)

func (s Application) String() string {
	var messages = map[Application]string{
		SystemError:      "システムエラー",
		StatusBadRequest: "リクエストの形式に誤りがあります",
		WrongPassword:    "パスワードが間違っています",
		LoginFailed:      "ログインに失敗しました",
		LogoutFailed:     "ログアウトに失敗しました",
		UnAuthorized:     "認証に失敗しました",
		Forbidden:        "この操作は許可されていません",
	}
	return messages[s]
}
