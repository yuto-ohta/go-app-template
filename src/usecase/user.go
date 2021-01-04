package usecase

import (
	"go-app-template/src/domain"
)

type UserUseCase interface {
	FindById(id int) (domain.User, error)
	CreateUser(userName string) (domain.User, error)
}
