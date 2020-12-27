package infrastructure

import (
	"go-app-template/config/db"
	"go-app-template/domain"
	"go-app-template/errors"
	"go-app-template/infrastructure/model"
	"gorm.io/gorm"
	"net/http"
)

type UserRepositoryImpl struct {
}

func NewUserRepositoryImpl() *UserRepositoryImpl {
	return &UserRepositoryImpl{}
}

func (u UserRepositoryImpl) FindById(id domain.UserId) (domain.User, error) {
	var userModel model.User
	var user domain.User
	var err error

	if err = db.Conn.Raw("SELECT * FROM users WHERE id = ?", id.GetValue()).Scan(&userModel).Error; err != nil {
		return user, err
	}

	user = *userModel.ToDomain()

	if user.GetId().GetValue() == 0 {
		err = errors.NewAppError(gorm.ErrRecordNotFound, http.StatusNotFound)
		return user, err
	}

	return user, nil
}
