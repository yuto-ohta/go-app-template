package impl

import (
	"errors"
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
	var userId *valueobject.UserId
	var err error
	if userId, err = valueobject.NewUserIdWithId(id); err != nil {
		return domain.User{}, apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}
	return u.userRepository.FindById(*userId)
}

func (u UserUseCaseImpl) CreateUser(userName string) (domain.User, error) {
	var (
		user   *domain.User
		err    error
		appErr *apperror.AppError
	)

	if user, err = domain.NewUser(userName); err != nil {
		if errors.Is(err, appErr) {
			return domain.User{}, apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
		}
		return domain.User{}, apperror.NewAppError(err)
	}

	return u.userRepository.CreateUser(*user)
}
