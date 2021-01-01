package repository

import (
	"go-app-template/src/domain"
	"go-app-template/src/domain/values"
)

type UserRepository interface {
	FindById(id values.UserId) (domain.User, error)
	CreateUser(user domain.User) (domain.User, error)
}
