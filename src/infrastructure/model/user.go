package model

import (
	"go-app-template/src/domain"
	"go-app-template/src/domain/value"
)

type User struct {
	ID   int
	Name string
}

func (u User) ToDomain() *domain.User {
	id := value.NewUserIdWithId(u.ID)
	name := u.Name
	return domain.NewUserWithUserId(*id, name)
}
