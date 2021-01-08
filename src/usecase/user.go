package usecase

import (
	"go-app-template/src/api/controller/dto"
)

type UserUseCase interface {
	FindById(id int) (dto.UserDto, error)
	FindAll() ([]dto.UserDto, error)
	CreateUser(userName string) (dto.UserDto, error)
	DeleteUser(id int) (dto.UserDto, error)
	UpdateUser(id int, user dto.UserDto) (dto.UserDto, error)
}
