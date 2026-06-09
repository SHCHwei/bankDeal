package model

import (
	"time"
	"database/sql"
)

type Deal struct {
	ID                  int			
    AccountID           int			`schema:"account_id" validate:"required"`
    Volume              int64		`schema:"volume" validate:"required,min=1,max=50000"`
    TransactionType     uint8		`schema:"transaction_type"`
    TradingAccountID    int			`schema:"trading_account_id"`
    Remark              string		`schema:"remark"`
	CreatedAt           time.Time
}

type DealRepository interface {
	FindAll() ([]*Deal, error)
	FindByID(id int) (*Deal, error)
	Save(tx *sql.Tx, deal *Deal) error
}

type DealService interface {
	ListDeals() ([]*Deal, error)
	GetDeal(id int) (map[string]string, error)
	CreateDeal(accountID int, volume int64, transactionType uint8, tradingAccountID int, remark string) (*Deal, error)
}
