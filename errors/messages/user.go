package messages

type User int

const (
	InvalidUserId User = iota
	UserNotFound
)

func (e User) String() string {
	switch e {
	case InvalidUserId:
		return "ユーザーIDの形式が間違っています"
	case UserNotFound:
		return "対象のユーザーが見つかりません"
	default:
		return "システムエラー"
	}
}
