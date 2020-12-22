package usecase

import (
	"go-app-template/domain"
)

type UserUseCase interface {
	FindById(id int) (domain.User, error)
}
