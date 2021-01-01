package value

type UserId struct {
	id int
}

func NewUserId(id int) *UserId {
	return &UserId{id: id}
}

func (u UserId) GetValue() int {
	return u.id
}
