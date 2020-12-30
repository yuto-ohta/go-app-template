package impl

import (
	"go-app-template/src/domain"
	"go-app-template/src/domain/repository"
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
