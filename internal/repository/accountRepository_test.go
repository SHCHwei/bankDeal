package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"bankDeal/internal/database"
	"bankDeal/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAccountRepository(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer mockDB.Close()
	oldDB := database.MariaDB
	defer func() { database.MariaDB = oldDB }()
	database.MariaDB = mockDB

	repo := NewAccountRepository(database.MariaDB)

	t.Run("FindByID returns account", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "userID", "bankID", "accountName", "balance", "createdAt", "updatedAt"}).
			AddRow(1, 42, 7, "savings", int64(1200), time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM accounts WHERE id = ?")).
			WithArgs(1).
			WillReturnRows(rows)

		account, err := repo.FindByID(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if account.ID != 1 || account.UserID != 42 || account.BankID != 7 || account.AccountName != "savings" || account.Balance != 1200 {
			t.Fatalf("unexpected account returned: %+v", account)
		}
	})

	t.Run("FindByUserID returns account", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "userID", "bankID", "accountName", "balance", "createdAt", "updatedAt"}).
			AddRow(5, 77, 9, "checking", int64(2500), time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM accounts WHERE user_id = ?")).
			WithArgs(77).
			WillReturnRows(rows)

		account, err := repo.FindByUserID(77)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if account.ID != 5 || account.UserID != 77 || account.BankID != 9 || account.AccountName != "checking" || account.Balance != 2500 {
			t.Fatalf("unexpected account returned: %+v", account)
		}
	})

	t.Run("Create inserts account", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO accounts (userID, bankID, accountName, balance) VALUES (?, ?, ?, ?)")).
			WithArgs(11, 22, "new account", int64(3000)).
			WillReturnResult(sqlmock.NewResult(101, 1))
		mock.ExpectCommit()


		txManager := database.NewTxManager(mockDB)
		account := &model.Account{UserID: 11, BankID: 22, AccountName: "new account", Balance: 3000}
		err := txManager.RunInTransaction(context.Background(), func(txCtx context.Context) error {
			return repo.Create(txCtx, account)
		})
		if err != nil {
			t.Fatalf("Create returned error: %v", err)
		}
		if account.ID != 101 {
			t.Fatalf("expected ID 101, got %d", account.ID)
		}

	})

	t.Run("Update updates account", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE accounts SET accountName = ? , balance = ? WHERE id = ?")).
			WithArgs("updated", int64(4000), 101).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()


		txManager := database.NewTxManager(mockDB)
		account := &model.Account{ID: 101, AccountName: "updated", Balance: 4000}
		err := txManager.RunInTransaction(context.Background(), func(txCtx context.Context) error {
			return repo.Update(txCtx, account)
		})
		if err != nil {
			t.Fatalf("Update returned error: %v", err)
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
