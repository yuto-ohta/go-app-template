package dto

import "go-app-template/src/apperror"

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
	if err := ValidateId(l.Id); err != nil {
		return apperror.NewAppError(err)
	}
	if err := ValidatePassword(l.Password); err != nil {
		return apperror.NewAppError(err)
	}
	return nil
}
