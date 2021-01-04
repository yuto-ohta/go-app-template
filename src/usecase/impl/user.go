package impl

import (
	"go-app-template/src/apperror"
	"go-app-template/src/domain"
	"go-app-template/src/domain/repository"
	"go-app-template/src/domain/valueobject"
	"net/http"
)

type UserUseCaseImpl struct {
	userRepository repository.UserRepository
}

func NewUserUseCaseImpl(userRepository repository.UserRepository) *UserUseCaseImpl {
	return &UserUseCaseImpl{userRepository: userRepository}
}

func (u UserUseCaseImpl) FindById(id int) (domain.User, error) {
	userId := valueobject.NewUserIdWithId(id)
	return u.userRepository.FindById(*userId)
}

func (u UserUseCaseImpl) CreateUser(userName string) (domain.User, error) {
	var user *domain.User
	var err error
	if user, err = domain.NewUser(userName); err != nil {
		return *new(domain.User), apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}
	return u.userRepository.CreateUser(*user)
}
