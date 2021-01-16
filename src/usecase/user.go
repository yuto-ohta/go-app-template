package usecase

import (
	"go-app-template/src/api/controller/dto"
	"go-app-template/src/usecase/appmodel"
)

type UserUseCase interface {
	GetUser(id int) (dto.UserDto, error)
	GetAllUser(condition appmodel.SearchCondition) (dto.UserPage, error)
	CreateUser(userDto dto.UserDto) (dto.UserDto, error)
	DeleteUser(id int) (dto.UserDto, error)
	UpdateUser(id int, user dto.UserDto) (dto.UserDto, error)
}
