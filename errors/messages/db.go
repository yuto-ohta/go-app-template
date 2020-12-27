package messages

type DB int

const (
	RecordNotFound DB = iota
)

func (e DB) String() string {
	switch e {
	case RecordNotFound:
		return "対象のレコードが見つかりません"
	default:
		return "システムエラー"
	}
}
