package repository

import (
	"testing"
	
	"bankDeal/internal/database"
	"bankDeal/internal/model"
)


func TestDealRepository(t *testing.T) {
	// 1. 初始化 repository
	dealRepo := NewDealRepositories()

	tx, _ := database.MariaDB.Begin()

	// 2. 測試 Save 方法
	deal := &model.Deal{
		AccountID:       	1,
		Volume:          	1000,
		TransactionType: 	1,
		TradingAccountID:   2,
		Remark:          	"Test deal",
	}

	// 2.1 測試減款交易
	if err := dealRepo.Save(tx, deal); err != nil {
		t.Fatalf("Failed to save deal(minus): %v", err)
	}


	// 2.2 測試加款交易
	deal.TransactionType = 0
	deal.Remark = "Test deal 2"

	if err := dealRepo.Save(tx, deal); err != nil {
		t.Fatalf("Failed to save deal(plus): %v", err)
	}


	// 2.3 測試加款交易，volume 為負數
	deal.TransactionType = 0
	deal.Volume = -500
	deal.Remark = "Test deal 3: Volume = -500"

	if err := dealRepo.Save(tx, deal); err != nil {
		t.Fatalf("Failed to save deal(plus): %v", err)
	}



}
