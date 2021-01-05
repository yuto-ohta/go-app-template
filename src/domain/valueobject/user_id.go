package valueobject

import (
	"go-app-template/src/apperror"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

const (
	_notAllocated = -1
)

type UserId struct {
	id int
}

func NewUserId() *UserId {
	return &UserId{id: _notAllocated}
}

func NewUserIdWithId(id int) (*UserId, error) {
	userId := &UserId{id: id}
	if err := userId.IsValidResErr(); err != nil {
		return nil, apperror.NewAppError(err)
	}
	return userId, nil
}

func (u UserId) GetValue() int {
	return u.id
}

func (u UserId) IsValidResErr() error {
	if u.isAllocated() {
		rules := "gte=1"
		err := validate.Var(u.id, rules)
		if err != nil {
			return apperror.NewAppError(err)
		}
	}

	return nil
}

func (u UserId) isAllocated() bool {
	return u.id != _notAllocated
}
