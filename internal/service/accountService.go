package service

import (
	"fmt"
	"bankDeal/internal/model"
)

type accountService struct {
	repo model.AccountRepository
}


func (a* accountService) GetAccount(id int) (*model.Account, error) {

	
	account, err := a.repo.FindByID(id)

	if err != nil {
		return nil, fmt.Errorf("account %d not found: %w", id, err)
	}

	return account, nil
}