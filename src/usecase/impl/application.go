package impl

import (
	"fmt"
	"go-app-template/src/api/controller/dto"
	"go-app-template/src/apperror"
	"go-app-template/src/domain"
	"go-app-template/src/domain/repository"
	"go-app-template/src/domain/valueobject"
	"go-app-template/src/usecase/appmodel"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type ApplicationUseCaseImpl struct {
	userRepository repository.UserRepository
}

func NewApplicationUseCaseImpl(userRepository repository.UserRepository) *ApplicationUseCaseImpl {
	return &ApplicationUseCaseImpl{userRepository: userRepository}
}

/**************************************
	ログイン
**************************************/

func (a ApplicationUseCaseImpl) Login(loginDto dto.LoginReceiveDto) (appmodel.SignedToken, error) {
	var err error

	// get userId from loginDto
	var userId *valueobject.UserId
	if userId, err = valueobject.NewUserIdWithId(loginDto.Id); err != nil {
		return "", apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
	}

	// get user
	var user domain.User
	if user, err = a.userRepository.FindById(*userId); err != nil {
		return "", apperror.NewAppError(err)
	}

	// check password
	userPassword := user.GetPassword()
	if err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(loginDto.Password)); err != nil {
		return "", apperror.NewAppErrorWithStatus(fmt.Errorf("パスワードが間違っています"), http.StatusUnauthorized)
	}

	// generate token
	var token appmodel.SignedToken
	if token, err = appmodel.GenerateSignedToken(user.GetId().GetValue()); err != nil {
		return "", apperror.NewAppError(err)
	}

	return token, nil
}
