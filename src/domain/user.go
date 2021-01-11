package domain

import (
	"fmt"
	"go-app-template/src/api/controller/dto"
	"go-app-template/src/apperror"
	"go-app-template/src/domain/valueobject"
	"go-app-template/src/usecase/appmodel"
	"sort"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

const (
	id userSortColumn = iota + 1
	userName
)

type User struct {
	id   valueobject.UserId
	name string
}

type userSortColumn int

func (u userSortColumn) string() string {
	var messages = map[userSortColumn]string{
		id:       "id",
		userName: "name",
	}
	return messages[u]
}

func NewUser(name string) (*User, error) {
	user := &User{id: *valueobject.NewUserId(), name: name}
	if err := user.Validate(); err != nil {
		return nil, apperror.NewAppError(err)
	}
	return user, nil
}

func NewUserWithUserId(id valueobject.UserId, name string) (*User, error) {
	user := &User{id: id, name: name}
	if err := user.Validate(); err != nil {
		return nil, apperror.NewAppError(err)
	}
	return user, nil
}

func (u User) GetId() valueobject.UserId {
	return u.id
}

func (u User) GetName() string {
	return u.name
}

func (u User) ToDto() *dto.UserDto {
	return &dto.UserDto{
		Id:   u.id.GetValue(),
		Name: u.name,
	}
}

func (u User) Validate() error {
	if err := u.validateName(); err != nil {
		return apperror.NewAppError(err)
	}
	return u.GetId().Validate()
}

func (u User) validateName() error {
	rules := "min=1,max=8"
	if err := validate.Var(u.name, rules); err != nil {
		return apperror.NewAppError(err)
	}
	return nil
}

func Sort(orderBy string, order appmodel.Order, target []User) error {
	var err error

	var sortColumn userSortColumn
	if sortColumn, err = getSortColumn(orderBy); err != nil {
		return apperror.NewAppError(err)
	}

	switch sortColumn {
	case userName:
		sortByUserName(order, target)
	default:
		sortById(order, target)
	}
	return nil
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
	case id.string():
		return id, nil
	case userName.string():
		return userName, nil
	default:
		return -1, apperror.NewAppError(fmt.Errorf("指定のColumnはソートに使用できるものではありません, param: %v", param))
	}
}
