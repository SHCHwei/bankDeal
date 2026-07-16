package model

import (
	"bankDeal/internal/dto/request"
	"time"
	"context"
)

type User struct {
	ID        int
	FirstName string    `fake:"{firstname}"`
	LastName  string    `fake:"{lastname}"`
	Email     string    `fake:"{email}"`
	Phone     string    `fake:"{number:11111111, 99999999}"`
	BirthDate string
	CreatedAt time.Time `fake:"-"`
	UpdatedAt time.Time `fake:"-"`
}

type UserService interface {
	GetUser(id int) (*User, error)
	CreateUser(ctx context.Context, requestData request.CreateUser) error
}

type UserRepository interface {
	FindUserByID(id int) (*User, error)
	Search(user User) (*User, error)
	InsertUser(ctx context.Context, user User) (int64, error)
}
