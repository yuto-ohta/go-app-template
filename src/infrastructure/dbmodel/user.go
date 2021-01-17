package dbmodel

import (
	"go-app-template/src/apperror"
	"go-app-template/src/domain"
	"go-app-template/src/domain/valueobject"
)

type User struct {
	ID       int
	Name     string
	Password string
}

/**************************************
	Constructor
**************************************/
func NewUserModel(user domain.User) User {
	return User{
		ID:       user.GetId().GetValue(),
		Name:     user.GetName(),
		Password: user.GetPassword(),
	}
}

/**************************************
	Conversion
**************************************/
func (u User) ToDomain() (*domain.User, error) {
	var err error

	var userId *valueobject.UserId
	if userId, err = valueobject.NewUserIdWithId(u.ID); err != nil {
		return nil, apperror.NewAppError(err)
	}

	user := domain.NewUserBuilder().Id(*userId).Name(u.Name).Password(u.Password).BuildWithoutValidateAndEncrypt()
	return user, nil
}
