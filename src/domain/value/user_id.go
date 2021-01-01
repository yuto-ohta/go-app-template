package value

type UserId struct {
	id int
}

const notAllocated = -1

func NewUserId() *UserId {
	return &UserId{id: notAllocated}
}

func NewUserIdWithId(id int) *UserId {
	return &UserId{id: id}
}

func (u UserId) GetValue() int {
	return u.id
}

func (u UserId) IsAllocated() bool {
	return u.id != notAllocated
}
