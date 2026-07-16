package repository

import (
	"bankDeal/internal/database"
	"bankDeal/internal/model"
	"regexp"
	"testing"
	"context"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestFindUserByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("open db err: %v", err)
	}
	defer db.Close()
	database.MariaDB = db

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "FirstName", "LastName", "Email", "Phone", "BirthDate", "CreatedAt", "UpdatedAt"}).
		AddRow(1, "John", "Doe", "john@example.com", "0912345678", "1990-01-01", now, now)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE id = ?")).WithArgs(1).WillReturnRows(rows)

	repo := NewUserRepository(database.MariaDB)
	u, err := repo.FindUserByID(1)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if u.ID != 1 || u.FirstName != "John" {
		t.Fatalf("unexpected user: %+v", u)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestSearch_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("open db err: %v", err)
	}
	defer db.Close()
	database.MariaDB = db

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "FirstName", "LastName", "Email", "Phone", "BirthDate", "CreatedAt", "UpdatedAt"}).
		AddRow(2, "Jane", "Smith", "jane@example.com", "0987654321", "", now, now)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, FirstName, LastName, Email, Phone, BirthDate, CreatedAt, UpdatedAt FROM users WHERE " + "FirstName = ?" + " LIMIT 1")).WithArgs("Jane").WillReturnRows(rows)

	repo := NewUserRepository(database.MariaDB)
	res, err := repo.Search(model.User{FirstName: "Jane"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res == nil || res.FirstName != "Jane" {
		t.Fatalf("unexpected search result: %+v", res)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("open db err: %v", err)
	}
	defer db.Close()
	database.MariaDB = db

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (FirstName, LastName, Email, Phone, BirthDate) VALUES (?, ?, ?, ?, ?)")).
		WithArgs("Alice", "W", "alice@example.com", "0911000111", "2000-02-02").WillReturnResult(sqlmock.NewResult(77, 1))
	mock.ExpectCommit()

	repo := NewUserRepository(db)

	txManager := database.NewTxManager(db)
	var id int64
	err = txManager.RunInTransaction(context.Background(), func(txCtx context.Context) error {
		var insertErr error
		id, insertErr = repo.InsertUser(txCtx, model.User{FirstName: "Alice", LastName: "W", Email: "alice@example.com", Phone: "0911000111", BirthDate: "2000-02-02"})
		return insertErr
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if id != 77 {
		t.Fatalf("expected id 77, got %d", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}
