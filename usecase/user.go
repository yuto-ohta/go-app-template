package usecase

import (
	"go-app-template/domain"
)

type UserUseCase interface {
	FindById(id domain.UserId) (domain.User, error)
}
