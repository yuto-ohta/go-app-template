package domain

import (
	"encoding/json"
)

type User struct {
	id   UserId
	name string
}

func NewUser(id UserId, name string) *User {
	return &User{id: id, name: name}
}

func (u User) GetId() UserId {
	return u.id
}

func (u User) GetName() string {
	return u.name
}

func (u User) MarshalJSON() ([]byte, error) {
	value, err := json.Marshal(&struct {
		Id   int
		Name string
	}{
		Id:   u.id.GetValue(),
		Name: u.name,
	})
	return value, err
}

func (u *User) UnmarshalJSON(b []byte) error {
	pseudo := &struct {
		Id   int
		Name string
	}{}
	if err := json.Unmarshal(b, pseudo); err != nil {
		return err
	}

	u.id = *NewUserId(pseudo.Id)
	u.name = pseudo.Name

	return nil
}
