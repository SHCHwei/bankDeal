package database

import (
	"database/sql"
	"log"
	"time"

	"bankDeal/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var MariaDB *sql.DB
var err error
func MariaDBConnect(cfg config.Config) bool {
	dsn := "user:user1@tcp(bankdeal-db-1:3306)/bankDeal?charset=utf8mb4&parseTime=True&loc=Local"

	// 1. Open a handle to the DB (doesn't connect yet)
	MariaDB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Panicf("sql Open error: %v", err)
	} else {
		log.Println("sql Open success")
	}



	// 2. Set connection pool settings (Recommended for production)
	MariaDB.SetConnMaxLifetime(time.Minute * 5)
	MariaDB.SetMaxOpenConns(25)
	MariaDB.SetMaxIdleConns(25)

	// 3. Ping the database to verify the actual connection
	err = MariaDB.Ping()

	if err != nil {
		log.Panicf("Error Content : %v", err)
	} else {
		log.Println("Database connection verified successfully")
	}

	return true

}

// BuildTables 使用 golang-migrate 來管理資料庫遷移
func BuildTables() {

	dsn := "root:root@tcp(bankdeal-db-1:3306)/bankDeal?charset=utf8mb4&parseTime=True&loc=Local"

	// 1. Open a handle to the DB (doesn't connect yet)
	MariaDB_root, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Panicf("sql Open error: %v", err)
	} else {
		log.Println("sql Open success")
	}

	defer MariaDB_root.Close()

	driver, err := mysql.WithInstance(MariaDB_root, &mysql.Config{})

	if err != nil {
		log.Panicf("Error creating MySQL driver: %v \n", err)
	} else {
		log.Println("MySQL driver created successfully")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)


	if err != nil {
		log.Panicf("Error creating migrate instance: %v \n", err)
	} else {
		log.Println("Migrate instance created successfully")
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Panicf("Error occurred while applying migrations: %v", err)
	}

}