package infrastructure

import (
	"go-app-template/config/db"
	"go-app-template/domain"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
}

func (u UserRepositoryImpl) FindById(id int) (domain.User, error) {
	var user domain.User
	var err error

	if err = db.Db.Raw("SELECT * FROM users WHERE id = ?", id).Scan(&user).Error; err != nil {
		return user, err
	}

	if user.ID == 0 {
		err = gorm.ErrRecordNotFound
		return user, err
	}

	return user, nil
}
