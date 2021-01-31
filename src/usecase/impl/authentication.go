package impl

import (
	"fmt"
	"go-app-template/src/api/controller/dto"
	"go-app-template/src/apperror"
	"go-app-template/src/usecase"
	"go-app-template/src/usecase/appmodel"
	"go-app-template/src/usecase/appmodel/session"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthenticationUseCaseImpl struct {
	sessionManager session.Manager
	userUseCase    usecase.UserUseCase
}

func NewAuthenticationUseCaseImpl(manager session.Manager, userUseCase usecase.UserUseCase) *AuthenticationUseCaseImpl {
	return &AuthenticationUseCaseImpl{
		sessionManager: manager,
		userUseCase:    userUseCase,
	}
}

/**************************************
	認証
**************************************/

func (a AuthenticationUseCaseImpl) Authenticate(c echo.Context, targetUserIdInt int) (dto.UserResDto, error) {
	var err error

	// get session
	sess := a.sessionManager.GetSession(c)

	// get signedToken from session
	var token appmodel.SignedToken
	if token = appmodel.GetTokenFromSession(sess); len(token) == 0 {
		return dto.UserResDto{}, apperror.NewAppErrorWithStatus(fmt.Errorf("token not found"), http.StatusUnauthorized)
	}

	// parse token
	var authInfo *appmodel.AuthInfo
	if authInfo, err = appmodel.ParseSignedToken(token); err != nil {
		return dto.UserResDto{}, apperror.NewAppErrorWithStatus(err, http.StatusUnauthorized)
	}

	// get sessionUser
	var sessionUser dto.UserResDto
	if sessionUser, err = a.userUseCase.GetUser(authInfo.GetUserId().GetValue()); err != nil {
		return dto.UserResDto{}, apperror.NewAppError(err)
	}

	// check sessionUser and target are the same
	if sessionUser.Id != targetUserIdInt {
		return dto.UserResDto{}, apperror.NewAppErrorWithStatus(fmt.Errorf("sessionUser is not the target user, sessionUser: %v, targetUserId: %v", sessionUser, targetUserIdInt), http.StatusUnauthorized)
	}

	return sessionUser, nil
}
