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

type UserReceiveDto struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserResDto struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

/**************************************
	Validation
**************************************/

func (u UserReceiveDto) Validate() error {
	if err := ValidateId(u.Id); err != nil {
		return apperror.NewAppError(err)
	}
	if err := ValidateName(u.Name); err != nil {
		return apperror.NewAppError(err)
	}
	if err := ValidatePassword(u.Password); err != nil {
		return apperror.NewAppError(err)
	}
	return nil
}

func ValidateId(id int) error {
	rules := "omitempty,gte=1"
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

	// 最大8文字
	if utf8.RuneCountInString(name) > 8 {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`userNameは最大8文字までです, userName: %v`, name), http.StatusBadRequest)
	}

	return nil
}

func ValidatePassword(password string) error {
	const PasswordAllowedStr = apputil.UpCaseAlphabet + apputil.DownCaseAlphabet + apputil.Number

	// 空文字はNG
	if password == "" {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`"password"が空文字になっています, password: %v`, password), http.StatusBadRequest)
	}

	// 半角・全角スペース, 改行を含む場合はNG
	if apputil.ContainsSpace(password) {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`"password"に半角・全角スペース, 改行コードが含まれています, password: %v`, password), http.StatusBadRequest)
	}

	// 英数大文字小文字を含む最小8文字
	if err := validate.Var(password, fmt.Sprintf("containsany=%v", apputil.DownCaseAlphabet)); err != nil {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`passwordに小文字のアルファベットを入れてください, password: %v`, password), http.StatusBadRequest)
	}
	if err := validate.Var(password, fmt.Sprintf("containsany=%v", apputil.UpCaseAlphabet)); err != nil {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`passwordに大文字のアルファベットを入れてください, password: %v`, password), http.StatusBadRequest)
	}
	if err := validate.Var(password, fmt.Sprintf("containsany=%v", apputil.Number)); err != nil {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`passwordに数字を入れてください, password: %v`, password), http.StatusBadRequest)
	}
	if err := validate.Var(password, "min=8"); err != nil {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`passwordは8文字以上にしてください, password: %v`, password), http.StatusBadRequest)
	}

	// 英数大文字小文字以外の文字を含む場合はNG
	if !apputil.ContainsAllowedStrOnly(password, PasswordAllowedStr) {
		return apperror.NewAppErrorWithStatus(fmt.Errorf(`passwordには英数大文字小文字以外を含めることはできません, password: %v`, password), http.StatusBadRequest)
	}

	return nil
}
