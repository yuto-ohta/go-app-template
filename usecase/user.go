package usecase

import (
	"go-app-template/domain"
	"go-app-template/domain/repository"
)

type UserUseCase struct {
	UserRepository repository.UserRepository
}

func (u UserUseCase) FindById(id int) (domain.User, error) {
	user, err := u.UserRepository.FindById(id)
	return user, err
}
