package impl

import (
	"fmt"
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

func (u UserUseCaseImpl) FindById(id valueobject.UserId) (domain.User, error) {
	return u.userRepository.FindById(id)
}

func (u UserUseCaseImpl) CreateUser(user domain.User) (domain.User, error) {
	if !user.IsValidForRegister() {
		return user, apperror.NewAppError(fmt.Errorf("未登録のユーザーにuserIdが割り当てられています, user: %v", user), http.StatusInternalServerError)
	}
	return u.userRepository.CreateUser(user)
}
