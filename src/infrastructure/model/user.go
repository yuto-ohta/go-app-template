package model

import (
	"go-app-template/src/apperror"
	"go-app-template/src/domain"
	"go-app-template/src/domain/valueobject"
)

type User struct {
	ID   int
	Name string
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

	if user, err = domain.NewUserWithUserId(*userId, u.Name); err != nil {
		return nil, apperror.NewAppError(err)
	}

	return user, nil
}
