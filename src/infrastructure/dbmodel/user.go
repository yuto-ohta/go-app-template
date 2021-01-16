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

func NewUserModel(user domain.User) User {
	return User{
		ID:       user.GetId().GetValue(),
		Name:     user.GetName(),
		Password: user.GetPassword(),
	}
}

func (u User) ToDomain() (*domain.User, error) {
	var (
		err    error
		userId *valueobject.UserId
		user   *domain.User
	)

	if userId, err = valueobject.NewUserIdWithId(u.ID); err != nil {
		return nil, apperror.NewAppError(err)
	}

	if user, err = domain.NewUserWithUserId(*userId, u.Name, u.Password); err != nil {
		return nil, apperror.NewAppError(err)
	}

	return user, nil
}
