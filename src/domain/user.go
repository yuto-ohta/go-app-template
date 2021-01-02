package domain

import (
	"encoding/json"
	"go-app-template/src/domain/valueobject"
)

type User struct {
	id   valueobject.UserId
	name string
}

type userJSON struct {
	Id   int
	Name string
}

func NewUser(name string) *User {
	return &User{id: *valueobject.NewUserId(), name: name}
}

func NewUserWithUserId(id valueobject.UserId, name string) *User {
	return &User{id: id, name: name}
}

func (u User) GetId() valueobject.UserId {
	return u.id
}

func (u User) GetName() string {
	return u.name
}

func (u User) MarshalJSON() ([]byte, error) {
	value, err := json.Marshal(&userJSON{
		Id:   u.id.GetValue(),
		Name: u.name,
	})
	return value, err
}

func (u *User) UnmarshalJSON(b []byte) error {
	var userJSON userJSON
	if err := json.Unmarshal(b, &userJSON); err != nil {
		return err
	}

	u.id = *valueobject.NewUserIdWithId(userJSON.Id)
	u.name = userJSON.Name

	return nil
}

func (u User) IsValidForRegister() bool {
	return !u.GetId().IsAllocated()
}
