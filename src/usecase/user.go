package usecase

import (
	"go-app-template/src/domain"
	"go-app-template/src/domain/values"
)

type UserUseCase interface {
	FindById(id values.UserId) (domain.User, error)
	CreateUser(user domain.User) (domain.User, error)
}
