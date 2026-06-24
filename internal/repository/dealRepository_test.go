package repository

import (
	"regexp"
	"testing"
	"time"

	"bankDeal/internal/database"
	"bankDeal/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDealRepository(t *testing.T) {
	// initialize mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("open db err: %v", err)
	}
	defer db.Close()
	database.MariaDB = db

	// initialize repository
	dealRepo := NewDealRepositories()

	// prepare tx
	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin tx err: %v", err)
	}

	deal := &model.Deal{
		AccountID:        1,
		Volume:           1000,
		TransactionType:  1,
		TradingAccountID: 2,
		Remark:           "Test deal",
	}

	// expect exec for each save call
	now := time.Now()
	_ = now
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO deals (accountID, volume, transactionType, tradingAccountID, remark) VALUES (?, ?, ?, ?, ?)")).
		WithArgs(deal.AccountID, deal.Volume, deal.TransactionType, deal.TradingAccountID, deal.Remark).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := dealRepo.Save(tx, deal); err != nil {
		t.Fatalf("Failed to save deal(minus): %v", err)
	}

	// second save
	deal.TransactionType = 0
	deal.Remark = "Test deal 2"
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO deals (accountID, volume, transactionType, tradingAccountID, remark) VALUES (?, ?, ?, ?, ?)")).
		WithArgs(deal.AccountID, deal.Volume, deal.TransactionType, deal.TradingAccountID, deal.Remark).
		WillReturnResult(sqlmock.NewResult(2, 1))

	if err := dealRepo.Save(tx, deal); err != nil {
		t.Fatalf("Failed to save deal(plus): %v", err)
	}

	// third save: negative volume
	deal.TransactionType = 0
	deal.Volume = -500
	deal.Remark = "Test deal 3: Volume = -500"
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO deals (accountID, volume, transactionType, tradingAccountID, remark) VALUES (?, ?, ?, ?, ?)")).
		WithArgs(deal.AccountID, deal.Volume, deal.TransactionType, deal.TradingAccountID, deal.Remark).
		WillReturnResult(sqlmock.NewResult(3, 1))

	if err := dealRepo.Save(tx, deal); err != nil {
		t.Fatalf("Failed to save deal(plus): %v", err)
	}

	// ensure expectations met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}
