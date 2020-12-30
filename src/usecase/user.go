package usecase

import (
	"go-app-template/src/domain"
)

type UserUseCase interface {
	FindById(id domain.UserId) (domain.User, error)
}
