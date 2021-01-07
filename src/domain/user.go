package domain

import (
	"encoding/json"
	"go-app-template/src/apperror"
	"go-app-template/src/domain/valueobject"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

type User struct {
	id   valueobject.UserId
	name string
}

type userJSON struct {
	Id   int
	Name string
}

func NewUser(name string) (*User, error) {
	user := &User{id: *valueobject.NewUserId(), name: name}
	if err := user.Validate(); err != nil {
		return nil, apperror.NewAppError(err)
	}
	return user, nil
}

func NewUserWithUserId(id valueobject.UserId, name string) (*User, error) {
	user := &User{id: id, name: name}
	if err := user.Validate(); err != nil {
		return nil, apperror.NewAppError(err)
	}
	return user, nil
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
	var err error

	var userJSON userJSON
	if err = json.Unmarshal(b, &userJSON); err != nil {
		return err
	}

	var id *valueobject.UserId
	if id, err = valueobject.NewUserIdWithId(userJSON.Id); err != nil {
		return err
	}

	u.id = *id
	u.name = userJSON.Name
	return nil
}

func (u User) Validate() error {
	if err := u.validateName(); err != nil {
		return apperror.NewAppError(err)
	}

	return u.GetId().Validate()
}

func (u User) validateName() error {
	rules := "min=1,max=8"
	if err := validate.Var(u.name, rules); err != nil {
		return apperror.NewAppError(err)
	}
	return nil
}
