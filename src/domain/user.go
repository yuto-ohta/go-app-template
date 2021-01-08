package domain

import (
	"go-app-template/src/api/controller/dto"
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

func (u User) ToDto() *dto.UserDto {
	return &dto.UserDto{
		Id:   u.id.GetValue(),
		Name: u.name,
	}
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
