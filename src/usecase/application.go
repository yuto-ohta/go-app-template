package usecase

import "go-app-template/src/api/controller/dto"

type ApplicationUseCase interface {
	Login(loginDto dto.LoginReceiveDto) (dto.LoginResDto, error)
}
