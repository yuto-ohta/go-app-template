package repository

import (
	"go-app-template/src/domain"
	"go-app-template/src/domain/valueobject"
)

type UserRepository interface {
	FindById(id valueobject.UserId) (domain.User, error)
	CreateUser(user domain.User) (domain.User, error)
}
