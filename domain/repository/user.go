package repository

import "go-app-template/domain"

type UserRepository interface {
	FindById(id int) (domain.User, error)
}
