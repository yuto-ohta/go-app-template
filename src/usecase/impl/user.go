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

func (u UserUseCaseImpl) FindById(id int) (dto.UserDto, error) {
	var err error

	//get userId
	var userId *valueobject.UserId
	if userId, err = valueobject.NewUserIdWithId(id); err != nil {
		return dto.UserDto{}, apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}

	// find user
	var found domain.User
	if found, err = u.userRepository.FindById(*userId); err != nil {
		return dto.UserDto{}, apperror.NewAppError(err)
	}

	return *found.ToDto(), nil
}

func (u UserUseCaseImpl) CreateUser(userName string) (dto.UserDto, error) {
	var err error

	// get user domain
	var newUser *domain.User
	if newUser, err = domain.NewUser(userName); err != nil {
		return dto.UserDto{}, apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}

	// create user
	var created domain.User
	if created, err = u.userRepository.CreateUser(*newUser); err != nil {
		return dto.UserDto{}, apperror.NewAppError(err)
	}

	return *created.ToDto(), nil
}

func (u UserUseCaseImpl) DeleteUser(id int) (dto.UserDto, error) {
	var err error

	// get userId
	var userId *valueobject.UserId
	if userId, err = valueobject.NewUserIdWithId(id); err != nil {
		return dto.UserDto{}, apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}

	// delete user
	var deleted domain.User
	if deleted, err = u.userRepository.DeleteUser(*userId); err != nil {
		return dto.UserDto{}, apperror.NewAppError(err)
	}

	return *deleted.ToDto(), nil
}

func (u UserUseCaseImpl) UpdateUser(id int, user dto.UserDto) (dto.UserDto, error) {
	var err error

	// get userId
	var userId *valueobject.UserId
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
