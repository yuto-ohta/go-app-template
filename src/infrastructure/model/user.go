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
	id := value.NewUserId(u.ID)
	name := u.Name
	return domain.NewUser(*id, name)
}
