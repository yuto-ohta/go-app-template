package dto

type UserPage struct {
	Users    []UserResDto `json:"users"`
	PageInfo PageInfo     `json:"pageInfo"`
}

type UserSortColumn int

const (
	ID UserSortColumn = iota + 1
	USERNAME
)

func (u UserSortColumn) String() string {
	var messages = map[UserSortColumn]string{
		ID:       "id",
		USERNAME: "name",
	}
	return messages[u]
}
