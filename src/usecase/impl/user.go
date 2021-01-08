package impl

import (
	"go-app-template/src/api/controller/dto"
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
	//get userId
	var userId *valueobject.UserId
	var err error
	if userId, err = valueobject.NewUserIdWithId(id); err != nil {
		return domain.User{}, apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}

	return u.userRepository.FindById(*userId)
}

func (u UserUseCaseImpl) CreateUser(userName string) (domain.User, error) {
	// get user domain
	var user *domain.User
	var err error
	if user, err = domain.NewUser(userName); err != nil {
		return domain.User{}, apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}

	return u.userRepository.CreateUser(*user)
}

func (u UserUseCaseImpl) DeleteUser(id int) (domain.User, error) {
	// get userId
	var userId *valueobject.UserId
	var err error
	if userId, err = valueobject.NewUserIdWithId(id); err != nil {
		return domain.User{}, apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}

	return u.userRepository.DeleteUser(*userId)
}

func (u UserUseCaseImpl) UpdateUser(id int, user dto.UserDto) (dto.UserDto, error) {
	// get userId
	var userId *valueobject.UserId
	var err error
	if userId, err = valueobject.NewUserIdWithId(id); err != nil {
		return dto.UserDto{}, apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}

	// 存在チェック
	if _, err = u.userRepository.FindById(*userId); err != nil {
		return dto.UserDto{}, apperror.NewAppError(err)
	}

	// get user domain
	var newUser *domain.User
	if newUser, err = domain.NewUser(user.Name); err != nil {
		return dto.UserDto{}, apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}

	// update user
	var updated domain.User
	if updated, err = u.userRepository.UpdateUser(*userId, *newUser); err != nil {
		return dto.UserDto{}, apperror.NewAppError(err)
	}

	return *updated.ToDto(), nil
}
