package model

import (
	"context"
	"time"
)

type Account struct {
	ID          int
	UserID      int
	BankID      int
	AccountName string
	Balance     int64
	CreatedAt   time.Time `fake:"-"`
	UpdatedAt   time.Time `fake:"-"`
}

type AccountRepository interface {
	FindByID(id int) (*Account, error)
	FindByUserID(userID int) (*Account, error)
	Create(ctx context.Context, data *Account) error
	Update(ctx context.Context, data *Account) error
}

type AccountService interface {
	GetAccount(id int) (*Account, error)
}
