package usecase

import (
	"go-app-template/src/api/controller/dto"
	"go-app-template/src/usecase/appmodel"
)

type UserUseCase interface {
	GetUser(id int) (dto.UserResDto, error)
	GetAllUser(condition appmodel.SearchCondition) (dto.UserPage, error)
	CreateUser(userDto dto.UserReceiveDto) (dto.UserResDto, error)
	DeleteUser(id int) (dto.UserResDto, error)
	UpdateUser(id int, user dto.UserReceiveDto) (dto.UserResDto, error)
}
