package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"bankDeal/internal/database"
	"bankDeal/internal/model"
)

type accountRepository struct {
	mu       sync.RWMutex
	accounts map[int]*model.Account
}

func NewAccountRepository() model.AccountRepository {
	return &accountRepository{
		accounts: make(map[int]*model.Account),
	}

}

func (r *accountRepository) FindByID(id int) (*model.Account, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	rows, err := database.MariaDB.Query("SELECT * FROM accounts WHERE id = ?", id)

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
	r.mu.RLock()
	defer r.mu.RUnlock()

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

func (r *accountRepository) Create(tx *sql.Tx, data *model.Account) error {
	const insertAccountSQL = "INSERT INTO accounts (userID, bankID, accountName, balance) VALUES (?, ?, ?, ?)"

	result, err := tx.Exec(insertAccountSQL, data.UserID, data.BankID, data.AccountName, data.Balance)
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

func (r *accountRepository) Update(tx *sql.Tx, data *model.Account) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	const updateSql = "UPDATE accounts SET accountName = ? , balance = ? WHERE id = ?"

	_, err := tx.ExecContext(context.TODO(), updateSql, data.AccountName, data.Balance, data.ID)

	if err != nil {
		return err
	}

	return nil

}
