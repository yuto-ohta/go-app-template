package valueobject

type UserId struct {
	id int
}

const _notAllocated = -1

func NewUserId() *UserId {
	return &UserId{id: _notAllocated}
}

func NewUserIdWithId(id int) *UserId {
	return &UserId{id: id}
}

func (u UserId) GetValue() int {
	return u.id
}

func (u UserId) IsAllocated() bool {
	return u.id != _notAllocated
}
