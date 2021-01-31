package usecase

import (
	"go-app-template/src/api/controller/dto"
	"go-app-template/src/usecase/appmodel"
)

type ApplicationUseCase interface {
	Login(loginDto dto.LoginReceiveDto) (appmodel.SignedToken, error)
}
