package model

import (
	"go-app-template/src/domain"
	"go-app-template/src/domain/valueobject"
)

type User struct {
	ID   int
	Name string
}

func (u User) ToDomain() *domain.User {
	id := valueobject.NewUserIdWithId(u.ID)
	name := u.Name
	return domain.NewUserWithUserId(*id, name)
}
