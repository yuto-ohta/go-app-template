package model

import "go-app-template/domain"

type User struct {
	ID   int
	Name string
}

func (u User) ToDomain() *domain.User {
	id := domain.NewUserId(u.ID)
	name := u.Name
	return domain.NewUser(*id, name)
}
