package dto

import (
	"fmt"
	"go-app-template/src/apperror"
	"go-app-template/src/apputil"
	"net/http"
)

type LoginReceiveDto struct {
	Id       int    `json:"id"`
	Password string `json:"password"`
}

type LoginResDto struct {
	Token string `json:"token"`
}

/**************************************
	Validation
**************************************/
func (l LoginReceiveDto) Validate() error {
	if err := validateLoginUserId(l.Id); err != nil {
		return apperror.NewAppError(err)
	}
	if err := validateLoginPassword(l.Password); err != nil {
		return apperror.NewAppError(err)
	}
	return nil
}

/**************************************
	private
**************************************/
func validateLoginUserId(id int) error {
	rules := "gte=1"
	if err := validate.Var(id, rules); err != nil {
		return apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}
	return nil
}

func validateLoginPassword(password string) error {
	// 空文字はNG
	if password == "" {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`"password"が空文字になっています, password: %v`, password), http.StatusBadRequest)
	}

	// 半角・全角スペース, 改行を含む場合はNG
	if apputil.ContainsSpace(password) {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`"password"に半角・全角スペース, 改行コードが含まれています, password: %v`, password), http.StatusBadRequest)
	}

	return nil
}
