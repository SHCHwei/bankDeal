package repository

import (
	"bankDeal/internal/database"
	"bankDeal/internal/model"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetBankByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	database.MariaDB = db

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "Code", "BankName", "CapitalAmount", "CreatedAt", "UpdatedAt"}).
		AddRow(1, "C1", "Bank A", 1000, now, now)

	mock.ExpectQuery(regexp.QuoteMeta("select * from banks where id = ?")).WithArgs(1).WillReturnRows(rows)

	repo := NewBankRepository()
	b, err := repo.GetBankByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if b == nil {
		t.Fatalf("expected bank, got nil")
	}
	if b.ID != 1 || b.Code != "C1" {
		t.Fatalf("unexpected bank data: %+v", b)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetBankByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("open db err: %v", err)
	}
	defer db.Close()
	database.MariaDB = db

	rows := sqlmock.NewRows([]string{"id", "Code", "BankName", "CapitalAmount", "CreatedAt", "UpdatedAt"})
	mock.ExpectQuery(regexp.QuoteMeta("select * from banks where id = ?")).WithArgs(999).WillReturnRows(rows)

	repo := NewBankRepository()
	_, err = repo.GetBankByID(999)
	if err == nil {
		t.Fatalf("expected error for not found, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateBank_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("open db err: %v", err)
	}
	defer db.Close()
	database.MariaDB = db

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO banks (Code, BankName, CapitalAmount) VALUES (?, ?, ?)")).
		WithArgs("C2", "Bank B", 2000).WillReturnResult(sqlmock.NewResult(42, 1))
	mock.ExpectCommit()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin tx err: %v", err)
	}

	repo := NewBankRepository()
	id, err := repo.CreateBank(tx, model.Bank{Code: "C2", BankName: "Bank B", CapitalAmount: 2000})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if id != 42 {
		t.Fatalf("expected id 42, got %d", id)
	}

	// commit to satisfy sqlmock expectation
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit tx err: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}
