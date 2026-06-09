package repository

import (
	"bankDeal/internal/database"
	"bankDeal/internal/model"
	"database/sql"
	"fmt"
	"sync"
)



type bankRepository struct {
	mu       sync.RWMutex
	bankList map[int]*model.Bank
}


func NewBankRepository() model.BankRepository {
	return &bankRepository{
		bankList: make(map[int]*model.Bank),
	}
}


func(r *bankRepository) GetBankByID(id int)(*model.Bank, error){

	rows, err := database.MariaDB.Query("select * from banks where id = ?", id)

	if err != nil {
		return nil, fmt.Errorf("bank id can not null")
	}

	defer rows.Close()

	var bankData model.Bank


	if !rows.Next() {
		return nil, fmt.Errorf("bank %d not found", id)
	}


	if err := rows.Scan(&bankData.ID, &bankData.Code, &bankData.BankName, &bankData.CapitalAmount, &bankData.CreatedAt, &bankData.UpdatedAt); err != nil {
		return nil, fmt.Errorf("bank id %d is not found. %v", id, err.Error())
	}

	return &bankData, nil
}



func (r *bankRepository) CreateBank(tx *sql.Tx, inputData model.Bank) (int64, error) {

	const insertUserSQL = "INSERT INTO banks (Code, BankName, CapitalAmount) VALUES (?, ?, ?)"

	result, err := tx.Exec(insertUserSQL, inputData.Code, inputData.BankName, inputData.CapitalAmount)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil

}