package usecase

import (
	"go-app-template/src/domain"
	"go-app-template/src/domain/valueobject"
)

type UserUseCase interface {
	FindById(id valueobject.UserId) (domain.User, error)
	CreateUser(user domain.User) (domain.User, error)
}
