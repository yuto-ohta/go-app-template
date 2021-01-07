package infrastructure

import (
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
		user      *domain.User
		err       error
	)

	// SQL実行
	if err = db.Conn.Raw("SELECT * FROM users WHERE id = ?", id.GetValue()).Scan(&userModel).Error; err != nil {
		return domain.User{}, apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
	}

	// NotFoundのエラーハンドリング
	if userModel.ID == 0 {
		return domain.User{}, apperror.NewAppErrorWithStatus(gorm.ErrRecordNotFound, http.StatusNotFound)
	}

	// domainに変換
	if user, err = userModel.ToDomain(); err != nil {
		return domain.User{}, apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
	}

	return *user, nil
}

func (u UserRepositoryImpl) CreateUser(user domain.User) (domain.User, error) {
	var err error
	userModel := model.User{Name: user.GetName()}

	// SQL実行
	result := db.Conn.Create(&userModel)
	if err = result.Error; err != nil {
		return user, apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
	}

	// domainに変換
	var createdUser *domain.User
	if createdUser, err = userModel.ToDomain(); err != nil {
		return user, apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
	}

	return *createdUser, nil
}

func (u UserRepositoryImpl) DeleteUser(id valueobject.UserId) (domain.User, error) {
	var err error

	// get user
	var user domain.User
	if user, err = u.FindById(id); err != nil {
		return domain.User{}, apperror.NewAppError(err)
	}

	// SQL実行
	result := db.Conn.Delete(&model.User{}, id.GetValue())
	if err = result.Error; err != nil {
		return domain.User{}, apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
	}

	return user, nil
}
