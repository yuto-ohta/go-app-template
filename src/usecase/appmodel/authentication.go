package appmodel

import (
	"fmt"
	"go-app-template/src/apperror"
	"go-app-template/src/config"
	"go-app-template/src/domain/valueobject"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
)

var _signatureEnvKey = config.GetConfig()["signature_key"].(string)

type SignedToken string

type AuthInfo struct {
	userId valueobject.UserId
	exp    int64
}

/**************************************
	Getter & Setter
**************************************/

func (a AuthInfo) GetUserId() valueobject.UserId {
	return a.userId
}

func (a AuthInfo) GetExp() int64 {
	return a.exp
}

/**************************************
	Main
**************************************/

func GenerateSignedToken(userId int) (SignedToken, error) {
	var err error

	// create token with claims
	token := jwt.New(jwt.SigningMethodHS256)
	expiresAt := time.Now().Add(time.Hour * 1).Unix()
	token.Claims = jwt.MapClaims{
		"userId": userId,
		"exp":    expiresAt,
	}

	// generate signed token
	var signedTokenStr string
	if signedTokenStr, err = token.SignedString([]byte(_signatureEnvKey)); err != nil {
		return "", apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
	}

	return SignedToken(signedTokenStr), nil
}

func GetTokenFromSession(sess *sessions.Session) SignedToken {
	token := sess.Values["token"]
	if token == nil {
		return ""
	}
	return SignedToken(token.(string))
}

func ParseSignedToken(signedToken SignedToken) (*AuthInfo, error) {
	var err error

	var token *jwt.Token
	if token, err = validateAndDecrypt(signedToken); err != nil {
		return nil, apperror.NewAppError(err)
	}

	if token == nil {
		return nil, apperror.NewAppError(fmt.Errorf("token not found in %v", signedToken))
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, apperror.NewAppError(fmt.Errorf("claims not found in %v", signedToken))
	}

	userIdFloat64, ok := claims["userId"].(float64)
	if !ok {
		return nil, apperror.NewAppError(fmt.Errorf("userId not found in %v", claims))
	}
	var userId *valueobject.UserId
	if userId, err = valueobject.NewUserIdWithId(int(userIdFloat64)); err != nil {
		return nil, apperror.NewAppError(fmt.Errorf("invalid userId value: %v", userIdFloat64))
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return nil, apperror.NewAppError(fmt.Errorf("exp not found in %v", claims))
	}
	exp := int64(expFloat)

	return newAuthInfo(*userId, exp), nil
}

/**************************************
	Constructor
**************************************/

func newAuthInfo(userId valueobject.UserId, exp int64) *AuthInfo {
	return &AuthInfo{
		userId: userId,
		exp:    exp,
	}
}

/**************************************
	private
**************************************/

func (a SignedToken) getStrValue() string {
	return string(a)
}

func validateAndDecrypt(signedToken SignedToken) (*jwt.Token, error) {
	// validation内容 tokenが正常であること、およびexpireが現在時間を過ぎていないことをチェックする
	token, err := jwt.Parse(signedToken.getStrValue(), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("unexpected signing method: %v", token.Method)
		}
		return []byte(_signatureEnvKey), nil
	})
	if err != nil {
		return nil, apperror.NewAppError(err)
	}
	return token, nil
}
