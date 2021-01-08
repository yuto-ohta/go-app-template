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
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (u UserDto) Validate() error {
	if err := ValidateId(u.Id); err != nil {
		return apperror.NewAppError(err)
	}
	if err := ValidateName(u.Name); err != nil {
		return apperror.NewAppError(err)
	}

	return nil
}

func ValidateId(id int) error {
	rules := "gte=1"
	if err := validate.Var(id, rules); err != nil {
		return apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}
	return nil
}

func ValidateName(name string) error {
	// 空文字はNG
	if name == "" {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`"userName"が空文字になっています, userName: %v`, name), http.StatusBadRequest)
	}

	// 半角・全角スペース, 改行を含む場合はNG
	if apputil.ContainsSpace(name) {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`"userName"に半角・全角スペース, 改行コードが含まれています, userName: %v`, name), http.StatusBadRequest)
	}

	// 9文字以上はNG
	if utf8.RuneCountInString(name) > 8 {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`userNameは最大8文字までです, userName: %v`, name), http.StatusBadRequest)
	}

	return nil
}
