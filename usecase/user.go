package usecase

import (
	"go-app-template/domain"
	"go-app-template/domain/repository"
	"go-app-template/infrastructure"
)

type UserUseCase struct {
}

func (u UserUseCase) FindById(id int) (domain.User, error) {
	var userRepository repository.UserRepository
	userRepository = infrastructure.UserRepositoryImpl{}
	user, err := userRepository.FindById(id)
	return user, err
}
