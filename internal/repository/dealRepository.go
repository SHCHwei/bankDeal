package repository

import (
	"context"
	"database/sql"

	"bankDeal/internal/database"
	"bankDeal/internal/model"
)

type dealRepositories struct {
	store map[string]*model.Deal
	sqlDB *sql.DB
}

func NewDealRepositories(sqlDB *sql.DB) model.DealRepository {
	return &dealRepositories{
		store: make(map[string]*model.Deal),
		sqlDB: sqlDB,
	}
}

func (r *dealRepositories) FindAll() ([]*model.Deal, error) {
	const dealSql = "select id, accountID, volume, transactionType, tradingAccountID, remark, createdAt from deals order by id asc"

	rows, err := r.sqlDB.Query(dealSql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deals := make([]*model.Deal, 0)
	for rows.Next() {
		var deal model.Deal
		if err := rows.Scan(&deal.ID, &deal.AccountID, &deal.Volume, &deal.TransactionType, &deal.TradingAccountID, &deal.Remark, &deal.CreatedAt); err != nil {
			return nil, err
		}
		deals = append(deals, &deal)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return deals, nil
}

func (r *dealRepositories) FindByID(id int) (*model.Deal, error) {

	const dealSql = "select id, accountID, volume, transactionType, tradingAccountID, remark, createdAt from deals where id = ?"
	var deal model.Deal

	row := r.sqlDB.QueryRow(dealSql, id)

	err := row.Scan(&deal.ID, &deal.AccountID, &deal.Volume, &deal.TransactionType, &deal.TradingAccountID, &deal.Remark, &deal.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &deal, nil
}

func (r *dealRepositories) Save(ctx context.Context, deal *model.Deal) error {

	const dealSql = "INSERT INTO deals (accountID, volume, transactionType, tradingAccountID, remark) VALUES (?, ?, ?, ?, ?)"

	tx := database.GetExecutor(ctx, r.sqlDB)

	_, err := tx.ExecContext(ctx, dealSql, deal.AccountID, deal.Volume, deal.TransactionType, deal.TradingAccountID, deal.Remark)

	if err != nil {
		return err
	}

	return nil
}
