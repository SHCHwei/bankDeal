package service

import (	
	"bankDeal/internal/model"
	"bankDeal/internal/database"
	"database/sql"

	"fmt"
)

type accountService struct {
	txManager   	*database.TxManager
	accountRepo 	model.AccountRepository
}


func NewAccountService(sqlDB *sql.DB, accountRepo model.AccountRepository) *accountService {
	return &accountService{
		txManager: database.NewTxManager(sqlDB),
		accountRepo: accountRepo,
	}
}


func (a* accountService) GetAccount(id int) (*model.Account, error) {

	
	account, err := a.accountRepo.FindByID(id)

	if err != nil {
		return nil, fmt.Errorf("account %d not found: %w", id, err)
	}

	return account, nil
}