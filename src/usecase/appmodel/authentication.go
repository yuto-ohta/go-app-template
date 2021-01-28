package appmodel

import (
	"go-app-template/src/apperror"
	"go-app-template/src/config"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var _signatureEnvKey = config.GetConfig()["signature_key"].(string)

type AuthToken string

func GenerateAuthToken(userId int, userName string) (AuthToken, error) {
	var err error

	// create token with claims
	token := jwt.New(jwt.SigningMethodHS256)
	expiresAt := time.Now().Add(time.Hour * 1).Unix()
	token.Claims = jwt.MapClaims{
		"sub":  userId,
		"name": userName,
		"exp":  expiresAt,
	}

	// generate signed token
	var signedTokenStr string
	if signedTokenStr, err = token.SignedString([]byte(_signatureEnvKey)); err != nil {
		return "", apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
	}

	return AuthToken(signedTokenStr), nil
}
