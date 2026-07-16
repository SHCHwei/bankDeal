package repository

import (
	"bankDeal/internal/database"
	"bankDeal/internal/model"
	"context"
	"database/sql"
	"fmt"
	
)

type accountRepository struct {
	accounts map[int]*model.Account
	sqlDB *sql.DB
}

func NewAccountRepository(sqlDB *sql.DB) model.AccountRepository {
	return &accountRepository{
		accounts: make(map[int]*model.Account),
		sqlDB: sqlDB,
	}

}

func (r *accountRepository) FindByID(id int) (*model.Account, error) {
	rows, err := r.sqlDB.Query("SELECT * FROM accounts WHERE id = ?", id)

	if err != nil {
		return nil, fmt.Errorf("account %d not found", id)
	}

	defer rows.Close()

	var targetAccount model.Account

	if !rows.Next() {
		return nil, fmt.Errorf("account %d not found", id)
	}

	if err := rows.Scan(&targetAccount.ID, &targetAccount.UserID, &targetAccount.BankID, &targetAccount.AccountName, &targetAccount.Balance, &targetAccount.CreatedAt, &targetAccount.UpdatedAt); err != nil {
		return nil, err
	}

	return &targetAccount, nil
}

func (r *accountRepository) FindByUserID(userID int) (*model.Account, error) {

	rows, err := database.MariaDB.Query("SELECT * FROM accounts WHERE user_id = ?", userID)

	if err != nil {
		return nil, fmt.Errorf("account for user %d not found", userID)
	}

	defer rows.Close()

	var targetAccount model.Account

	if !rows.Next() {
		return nil, fmt.Errorf("account for user %d not found", userID)
	}

	if err := rows.Scan(&targetAccount.ID, &targetAccount.UserID, &targetAccount.BankID, &targetAccount.AccountName, &targetAccount.Balance, &targetAccount.CreatedAt, &targetAccount.UpdatedAt); err != nil {
		return nil, err
	}

	return &targetAccount, nil
}

func (r *accountRepository) Create(ctx context.Context, data *model.Account) error {

	const insertAccountSQL = "INSERT INTO accounts (userID, bankID, accountName, balance) VALUES (?, ?, ?, ?)"
	tx := database.GetExecutor(ctx, r.sqlDB)
	result, err := tx.ExecContext(ctx, insertAccountSQL, data.UserID, data.BankID, data.AccountName, data.Balance)
	
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	data.ID = int(id)

	return nil
}

func (r *accountRepository) Update(ctx context.Context, data *model.Account) error {

	const updateSql = "UPDATE accounts SET accountName = ? , balance = ? WHERE id = ?"

	tx := database.GetExecutor(ctx, r.sqlDB)

	_, err := tx.ExecContext(ctx, updateSql, data.AccountName, data.Balance, data.ID)

	if err != nil {
		return err
	}

	return nil

}
