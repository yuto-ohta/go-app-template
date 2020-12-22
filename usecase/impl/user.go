package impl

import (
	"go-app-template/domain"
	"go-app-template/domain/repository"
)

type UserUseCaseImpl struct {
	UserRepository repository.UserRepository
}

func (u UserUseCaseImpl) FindById(id int) (domain.User, error) {
	user, err := u.UserRepository.FindById(id)
	return user, err
}
