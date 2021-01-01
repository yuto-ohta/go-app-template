package repository

import (
	"go-app-template/src/domain"
	"go-app-template/src/domain/value"
)

type UserRepository interface {
	FindById(id value.UserId) (domain.User, error)
	CreateUser(user domain.User) (domain.User, error)
}
