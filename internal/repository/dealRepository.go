package repository

import (
	"context"
	"database/sql"
	"sync"

	"bankDeal/internal/database"
	"bankDeal/internal/model"
)

type dealRepositories struct {
	mu    sync.RWMutex
	store map[string]*model.Deal
}

func NewDealRepositories() model.DealRepository {
	return &dealRepositories{
		store: make(map[string]*model.Deal),
	}
}

func (r *dealRepositories) FindAll() ([]*model.Deal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	deals := make([]*model.Deal, 0, len(r.store))
	for _, d := range r.store {
		deals = append(deals, d)
	}
	return deals, nil
}

func (r *dealRepositories) FindByID(id int) (*model.Deal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	const dealSql = "select id, accountID, volume, transactionType, tradingAccountID, remark, createdAt from deals where id = ?"

	var deal model.Deal

	row := database.MariaDB.QueryRow(dealSql, id)

	err := row.Scan(&deal.ID, &deal.AccountID, &deal.Volume, &deal.TransactionType, &deal.TradingAccountID, &deal.Remark, &deal.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &deal, nil
}

func (r *dealRepositories) Save(tx *sql.Tx, deal *model.Deal) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	const dealSql = "INSERT INTO deals (accountID, volume, transactionType, tradingAccountID, remark) VALUES (?, ?, ?, ?, ?)"

	_, err := tx.ExecContext(context.TODO(), dealSql, deal.AccountID, deal.Volume, deal.TransactionType, deal.TradingAccountID, deal.Remark)

	if err != nil {
		return err
	}

	return nil
}
