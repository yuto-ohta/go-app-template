package message

type User int

const (
	InvalidUserId User = iota + 1
	UserNotFound
	InvalidUserName
	CreateUserFailed
)

func (u User) String() string {
	var messages = map[User]string{
		InvalidUserId:    "ユーザーIDの形式が間違っています",
		UserNotFound:     "対象のユーザーが見つかりません",
		InvalidUserName:  "ユーザー名の形式が間違っています",
		CreateUserFailed: "ユーザー登録に失敗しました",
	}
	return messages[u]
}
