package repository

import "go-app-template/domain"

type UserRepository interface {
	FindById(id domain.UserId) (domain.User, error)
}
