package impl

import (
	"go-app-template/domain"
	"go-app-template/domain/repository"
)

type UserUseCaseImpl struct {
	userRepository repository.UserRepository
}

func NewUserUseCaseImpl(userRepository repository.UserRepository) *UserUseCaseImpl {
	return &UserUseCaseImpl{userRepository: userRepository}
}

func (u UserUseCaseImpl) FindById(id domain.UserId) (domain.User, error) {
	return u.userRepository.FindById(id)
}
