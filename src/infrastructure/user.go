package infrastructure

import (
	"fmt"
	"go-app-template/src/apperror"
	"go-app-template/src/config/db"
	"go-app-template/src/domain"
	"go-app-template/src/domain/valueobject"
	"go-app-template/src/infrastructure/model"
	"net/http"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
}

func NewUserRepositoryImpl() *UserRepositoryImpl {
	return &UserRepositoryImpl{}
}

func (u UserRepositoryImpl) FindById(id valueobject.UserId) (domain.User, error) {
	var (
		userModel model.User
		user      domain.User
		err       error
	)

	if err = db.Conn.Raw("SELECT * FROM users WHERE id = ?", id.GetValue()).Scan(&userModel).Error; err != nil {
		return user, err
	}

	user = *userModel.ToDomain()

	if user.GetId().GetValue() == 0 {
		err = apperror.NewAppErrorWithStatus(gorm.ErrRecordNotFound, http.StatusNotFound)
		return user, err
	}

	return user, nil
}

func (u UserRepositoryImpl) CreateUser(user domain.User) (domain.User, error) {
	userModel := model.User{Name: user.GetName()}
	result := db.Conn.Create(&userModel)

	if err := result.Error; err != nil {
		return user, apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
	}
	if rowsAffected := result.RowsAffected; rowsAffected != 1 {
		return user, apperror.NewAppErrorWithStatus(fmt.Errorf("INSERT文のRowsAffectedが1以外になっています, RowsAffected: %v", rowsAffected), http.StatusInternalServerError)
	}

	createdUser := *userModel.ToDomain()
	return createdUser, nil
}
