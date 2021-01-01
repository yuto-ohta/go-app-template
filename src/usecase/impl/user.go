package impl

import (
	"go-app-template/src/domain"
	"go-app-template/src/domain/repository"
	"go-app-template/src/domain/value"
)

type UserUseCaseImpl struct {
	userRepository repository.UserRepository
}

func NewUserUseCaseImpl(userRepository repository.UserRepository) *UserUseCaseImpl {
	return &UserUseCaseImpl{userRepository: userRepository}
}

func (u UserUseCaseImpl) FindById(id value.UserId) (domain.User, error) {
	return u.userRepository.FindById(id)
}
