package usecase

import (
	"go-app-template/src/api/controller/dto"
)

type UserUseCase interface {
	GetUser(id int) (dto.UserDto, error)
	GetAllUser(limit int, offset int) ([]dto.UserDto, error)
	CreateUser(userName string) (dto.UserDto, error)
	DeleteUser(id int) (dto.UserDto, error)
	UpdateUser(id int, user dto.UserDto) (dto.UserDto, error)
}
