package service

import (
	"bankDeal/internal/model"
	"bankDeal/internal/database"
	"database/sql"
	
	"fmt"
	"strconv"
	"context"

)

type dealService struct {
	txManager   	*database.TxManager
	dealRepo 		model.DealRepository
	accountRepo 	model.AccountRepository
}

func NewDealService(sqlDB *sql.DB, dealRepo model.DealRepository, accountRepo model.AccountRepository) model.DealService {
	return &dealService{
		txManager: database.NewTxManager(sqlDB),
		dealRepo: dealRepo, 
		accountRepo: accountRepo,
	}
}

func (s *dealService) ListDeals() ([]*model.Deal, error) {
	return s.dealRepo.FindAll()
}

func (s *dealService) GetDeal(id int) (map[string]string, error) {


	deal, err := s.dealRepo.FindByID(id)

	if err != nil {
		return nil , err
	}
	
	tradingAccount, err := s.accountRepo.FindByID(deal.TradingAccountID)

	if err != nil {
		return nil , err
	}

	dealDetail := make(map[string]string)


	dealDetail["DealAccount_id"] = strconv.Itoa(tradingAccount.ID)
	dealDetail["DealAccountName"] = tradingAccount.AccountName
	dealDetail["Volume"] = strconv.FormatInt(deal.Volume, 10)
	dealDetail["Remark"] = deal.Remark

	return dealDetail, nil
}

func (s *dealService) CreateDeal(ctx context.Context, accountID int, volume int64, transactionType uint8, tradingAccountID int, remark string) (*model.Deal, error) {


	var account , tradingAccount *model.Account
	var err error

	// 檢查甲方帳戶內狀態是否允許交易
	if accountID > 0 {
		
		account, err = s.accountRepo.FindByID(accountID)

		if err != nil {
			return nil, fmt.Errorf("your account %d not found : %s", accountID , err.Error())
		}

		if account.Balance < volume && transactionType == 1 {
			return nil, fmt.Errorf("your account %d has insufficient balance", accountID)
		}

	} else {
		return nil, fmt.Errorf("accountID is required")
	}


	// 檢查乙方帳戶內狀態是否允許交易
	if tradingAccountID > 0 {

		tradingAccount, err = s.accountRepo.FindByID(tradingAccountID)

		if err != nil {
			return nil, fmt.Errorf("trading account %d not found", tradingAccountID)
		}

		if tradingAccount.Balance < volume && transactionType == 0 {
			return nil, fmt.Errorf("trading account %d has insufficient balance", tradingAccountID)
		}

	} else {
		return nil, fmt.Errorf("tradingAccountID is required")
	}


	// 建立交易紀錄
	deal := &model.Deal{
		AccountID:       	accountID,
		Volume:          	volume,
		TransactionType: 	transactionType,
		TradingAccountID:   tradingAccountID,
		Remark:          	remark,
	}

	err = s.txManager.RunInTransaction(ctx, func(txCtx context.Context) error {

		if err := s.dealRepo.Save(txCtx, deal); err != nil {
			return err
		}


		// 帳戶異動紀錄 (甲方)
		if transactionType == 0 {
			account.Balance = account.Balance + volume
		} else {
			account.Balance = account.Balance - volume
		}

		if err := s.accountRepo.Update(txCtx, account); err != nil {
			return err
		}


		// 帳戶異動紀錄 (乙方)
		if transactionType == 1 {
			tradingAccount.Balance = tradingAccount.Balance + volume
		} else {
			tradingAccount.Balance = tradingAccount.Balance - volume
		}

		if err := s.accountRepo.Update(txCtx, tradingAccount); err != nil {
			return err
		}
	
		return nil
	})


	if err != nil {
		return nil, err
	}
	
	return deal, nil
}
