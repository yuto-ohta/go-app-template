package messages

type User int

const (
	InvalidUserId User = iota
	UserNotFound
)

func (u User) String() string {
	var messages = map[User]string{
		InvalidUserId: "ユーザーIDの形式が間違っています",
		UserNotFound:  "対象のユーザーが見つかりません",
	}
	return messages[u]
}
