package dto

import (
	"fmt"
	"go-app-template/src/apperror"
	"go-app-template/src/apputil"
	"net/http"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type UserDto struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

/**************************************
	Validation
**************************************/

func (u UserDto) Validate() error {
	if err := validateId(u.Id); err != nil {
		return apperror.NewAppError(err)
	}
	if err := validateName(u.Name); err != nil {
		return apperror.NewAppError(err)
	}
	if err := validatePassword(u.Password); err != nil {
		return apperror.NewAppError(err)
	}
	return nil
}

func validateId(id int) error {
	rules := "omitempty,gte=1"
	if err := validate.Var(id, rules); err != nil {
		return apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}
	return nil
}

func validateName(name string) error {
	// 空文字はNG
	if name == "" {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`"userName"が空文字になっています, userName: %v`, name), http.StatusBadRequest)
	}

	// 半角・全角スペース, 改行を含む場合はNG
	if apputil.ContainsSpace(name) {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`"userName"に半角・全角スペース, 改行コードが含まれています, userName: %v`, name), http.StatusBadRequest)
	}

	// 最大8文字
	if utf8.RuneCountInString(name) > 8 {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`userNameは最大8文字までです, userName: %v`, name), http.StatusBadRequest)
	}

	return nil
}

func validatePassword(password string) error {
	// 空文字はNG
	if password == "" {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`"password"が空文字になっています, password: %v`, password), http.StatusBadRequest)
	}

	// 半角・全角スペース, 改行を含む場合はNG
	if apputil.ContainsSpace(password) {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`"password"に半角・全角スペース, 改行コードが含まれています, password: %v`, password), http.StatusBadRequest)
	}

	// 英数大文字小文字を含む最小8文字
	if err := validate.Var(password, "containsany=abcdefghijklmnopqrstuvwsyz"); err != nil {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`passwordに小文字のアルファベットを入れてください, password: %v`, password), http.StatusBadRequest)
	}
	if err := validate.Var(password, "containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ"); err != nil {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`passwordに大文字のアルファベットを入れてください, password: %v`, password), http.StatusBadRequest)
	}
	if err := validate.Var(password, "containsany=0123456789"); err != nil {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`passwordに数字を入れてください, password: %v`, password), http.StatusBadRequest)
	}
	if err := validate.Var(password, "min=8"); err != nil {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`passwordは8文字以上にしてください, password: %v`, password), http.StatusBadRequest)
	}
	return nil
}
