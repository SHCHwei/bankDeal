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

func TestDealRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("open db err: %v", err)
	}
	defer db.Close()
	database.MariaDB = db

	repo := NewDealRepositories(db)

	t.Run("FindAll returns all deals", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "accountID", "volume", "transactionType", "tradingAccountID", "remark", "createdAt"}).
			AddRow(1, 10, int64(1000), uint8(1), 20, "first deal", now).
			AddRow(2, 11, int64(2000), uint8(2), 21, "second deal", now.Add(time.Second))

		mock.ExpectQuery(
			regexp.QuoteMeta("select id, accountID, volume, transactionType, tradingAccountID, remark, createdAt from deals order by id asc"),
		).WillReturnRows(rows)

		deals, err := repo.FindAll()
		if err != nil {
			t.Fatalf("FindAll returned error: %v", err)
		}
		if len(deals) != 2 {
			t.Fatalf("expected 2 deals, got %d", len(deals))
		}
		if deals[0].Remark != "first deal" || deals[1].Remark != "second deal" {
			t.Fatalf("unexpected deal rows returned: %+v", deals)
		}
	})

	t.Run("FindByID returns a deal", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "accountID", "volume", "transactionType", "tradingAccountID", "remark", "createdAt"}).
			AddRow(7, 12, int64(3000), uint8(1), 22, "specific deal", now)

		mock.ExpectQuery(
			regexp.QuoteMeta("select id, accountID, volume, transactionType, tradingAccountID, remark, createdAt from deals where id = ?"),
		).WithArgs(7).WillReturnRows(rows)

		deal, err := repo.FindByID(7)
		if err != nil {
			t.Fatalf("FindByID returned error: %v", err)
		}
		if deal == nil {
			t.Fatal("expected deal, got nil")
		}
		if deal.ID != 7 || deal.Remark != "specific deal" || deal.Volume != 3000 {
			t.Fatalf("unexpected deal returned: %+v", deal)
		}
	})

	t.Run("Save inserts a deal", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO deals (accountID, volume, transactionType, tradingAccountID, remark) VALUES (?, ?, ?, ?, ?)")).WithArgs(1, int64(1000), uint8(1), 2, "Test deal").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		deal := &model.Deal{
			AccountID:        1,
			Volume:           1000,
			TransactionType:  1,
			TradingAccountID: 2,
			Remark:           "Test deal",
		}

		txManager := database.NewTxManager(db)
		err = txManager.RunInTransaction(context.Background(), func(txCtx context.Context) error {
			return repo.Save(txCtx, deal)
		})
		if err != nil {
			t.Fatalf("Failed to save deal: %v", err)
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}
