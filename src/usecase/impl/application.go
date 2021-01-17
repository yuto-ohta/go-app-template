package impl

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go-app-template/src/api/controller/dto"
	"go-app-template/src/apperror"
	"go-app-template/src/domain"
	"go-app-template/src/domain/repository"
	"go-app-template/src/domain/valueobject"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

const _signatureEnvKey = "SIGNATURE_KEY"

type ApplicationUseCaseImpl struct {
	userRepository repository.UserRepository
}

func NewApplicationUseCaseImpl(userRepository repository.UserRepository) *ApplicationUseCaseImpl {
	return &ApplicationUseCaseImpl{userRepository: userRepository}
}

/**************************************
	ログイン
**************************************/
func (a ApplicationUseCaseImpl) Login(loginDto dto.LoginReceiveDto) (dto.LoginResDto, error) {
	var err error

	// get userId from loginDto
	var userId *valueobject.UserId
	if userId, err = valueobject.NewUserIdWithId(loginDto.Id); err != nil {
		return dto.LoginResDto{}, apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
	}

	// get user
	var user domain.User
	if user, err = a.userRepository.FindById(*userId); err != nil {
		return dto.LoginResDto{}, apperror.NewAppError(err)
	}

	// check password
	userPassword := user.GetPassword()
	if err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(loginDto.Password)); err != nil {
		return dto.LoginResDto{}, apperror.NewAppErrorWithStatus(fmt.Errorf("パスワードが間違っています"), http.StatusUnauthorized)
	}

	// create token with claims
	token := jwt.New(jwt.SigningMethodHS256)
	expiresAt := time.Now().Add(time.Hour * 1).Unix()
	token.Claims = jwt.MapClaims{
		"sub":  user.GetId().GetValue(),
		"name": user.GetName(),
		"exp":  expiresAt,
	}

	// generate signed token
	var signedTokenStr string
	var signedToken dto.LoginResDto
	if signedTokenStr, err = token.SignedString([]byte(os.Getenv(_signatureEnvKey))); err != nil {
		return dto.LoginResDto{}, apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
	}
	signedToken = dto.LoginResDto{Token: signedTokenStr}

	return signedToken, nil
}
