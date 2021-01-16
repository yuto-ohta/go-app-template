package domain

import (
	"fmt"
	"go-app-template/src/api/controller/dto"
	"go-app-template/src/apperror"
	"go-app-template/src/apputil"
	"go-app-template/src/domain/valueobject"
	"go-app-template/src/usecase/appmodel"
	"sort"

	"github.com/go-playground/validator/v10"
)

var _validate = validator.New()

const (
	_nameNotAllocated     = ""
	_passwordNotAllocated = ""
)

type User struct {
	id       valueobject.UserId
	name     string
	password string
}

/**************************************
	Builder
**************************************/

type UserBuilder struct {
	user *User
}

func NewUserBuilder() *UserBuilder {
	user := &User{
		id:       *valueobject.NewUserId(),
		name:     _nameNotAllocated,
		password: _passwordNotAllocated,
	}
	b := &UserBuilder{user: user}
	return b
}

func (b *UserBuilder) Id(id valueobject.UserId) *UserBuilder {
	b.user.id = id
	return b
}

func (b *UserBuilder) Name(name string) *UserBuilder {
	b.user.name = name
	return b
}

func (b *UserBuilder) Password(password string) *UserBuilder {
	b.user.password = password
	return b
}

func (b *UserBuilder) Build() (*User, error) {
	user := b.user
	if err := user.id.Validate(); err != nil {
		return nil, apperror.NewAppError(err)
	}
	if user.isNameAllocated() {
		if err := user.ValidateName(); err != nil {
			return nil, apperror.NewAppError(err)
		}
	}
	if user.isPasswordAllocated() {
		if err := user.ValidatePassword(); err != nil {
			return nil, apperror.NewAppError(err)
		}
		// hash化して詰め直す
		if err := encryptUserPassword(b.user); err != nil {
			return nil, apperror.NewAppError(err)
		}
	}
	return b.user, nil
}

func encryptUserPassword(user *User) error {
	hashedStr, err := apputil.GenerateHash(user.password)
	if err != nil {
		return apperror.NewAppError(err)
	}
	user.password = hashedStr
	return nil
}

/**************************************
	Getter & Setter
**************************************/

func (u User) GetId() valueobject.UserId {
	return u.id
}

func (u User) GetName() string {
	return u.name
}

func (u User) GetPassword() string {
	return u.password
}

/**************************************
	Conversion
**************************************/

func (u User) ToDto() *dto.UserResDto {
	return &dto.UserResDto{
		Id:   u.id.GetValue(),
		Name: u.name,
	}
}

/**************************************
	Validation
**************************************/

func (u User) Validate() error {
	if err := u.ValidateName(); err != nil {
		return apperror.NewAppError(err)
	}
	if err := u.ValidatePassword(); err != nil {
		return apperror.NewAppError(err)
	}
	return u.GetId().Validate()
}

func (u User) ValidateName() error {
	rules := "min=1,max=8"
	if err := _validate.Var(u.name, rules); err != nil {
		return apperror.NewAppError(err)
	}
	return nil
}

func (u User) ValidatePassword() error {
	if err := _validate.Var(u.name, "containsany=abcdefghijklmnopqrstuvwsyz"); err != nil {
		return apperror.NewAppError(fmt.Errorf(`passwordに小文字のアルファベットを入れてください, password: %v`, u.name))
	}
	if err := _validate.Var(u.name, "containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ"); err != nil {
		return apperror.NewAppError(fmt.Errorf(`passwordに大文字のアルファベットを入れてください, password: %v`, u.name))
	}
	if err := _validate.Var(u.name, "containsany=0123456789"); err != nil {
		return apperror.NewAppError(fmt.Errorf(`passwordに数字を入れてください, password: %v`, u.name))
	}
	if err := _validate.Var(u.name, "min=8"); err != nil {
		return apperror.NewAppError(fmt.Errorf(`passwordは8文字以上にしてください, password: %v`, u.name))
	}
	return nil
}

func (u User) isNameAllocated() bool {
	return u.name != _nameNotAllocated
}

func (u User) isPasswordAllocated() bool {
	return u.name != _passwordNotAllocated
}

/**************************************
	Sort
**************************************/

const (
	_id userSortColumn = iota + 1
	_userName
)

type userSortColumn int

func Sort(orderBy string, order appmodel.Order, target []User) error {
	var err error

	var sortColumn userSortColumn
	if sortColumn, err = getSortColumn(orderBy); err != nil {
		return apperror.NewAppError(err)
	}

	switch sortColumn {
	case _userName:
		sortByUserName(order, target)
	default:
		sortById(order, target)
	}
	return nil
}

func (u userSortColumn) string() string {
	var messages = map[userSortColumn]string{
		_id:       "id",
		_userName: "name",
	}
	return messages[u]
}

func sortById(order appmodel.Order, target []User) {
	if order == appmodel.ASC {
		sort.Slice(target, func(i, j int) bool {
			return target[i].id.GetValue() < target[j].id.GetValue()
		})
	} else {
		sort.Slice(target, func(i, j int) bool {
			return target[i].id.GetValue() > target[j].id.GetValue()
		})
	}
}

func sortByUserName(order appmodel.Order, target []User) {
	if order == appmodel.ASC {
		sort.Slice(target, func(i, j int) bool {
			return target[i].name < target[j].name
		})
	} else {
		sort.Slice(target, func(i, j int) bool {
			return target[i].name > target[j].name
		})
	}
}

func getSortColumn(param string) (userSortColumn, error) {
	switch param {
	case _id.string():
		return _id, nil
	case _userName.string():
		return _userName, nil
	default:
		return -1, apperror.NewAppError(fmt.Errorf("指定のColumnはソートに使用できるものではありません, param: %v", param))
	}
}
