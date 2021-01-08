package usecase

import (
	"go-app-template/src/api/controller/dto"
	"go-app-template/src/domain"
)

type UserUseCase interface {
	FindById(id int) (domain.User, error)
	CreateUser(userName string) (domain.User, error)
	DeleteUser(id int) (domain.User, error)
	UpdateUser(id int, user dto.UserDto) (dto.UserDto, error)
}
