package service

import (
	"errors"

	"bankDeal/internal/pkg"
	"bankDeal/internal/database"
	"bankDeal/internal/model"
	"bankDeal/internal/dto/request"
)

type userService struct {
	userRepo    model.UserRepository
	accountRepo model.AccountRepository
	bankRepo 	model.BankRepository
}

func NewUserService(userRepo model.UserRepository, accountRepo model.AccountRepository, bankRepo model.BankRepository) model.UserService {
	return &userService{
		userRepo: userRepo, 
		accountRepo: accountRepo,
		bankRepo: bankRepo,
	}
}

func (s *userService) GetUser(id int) (*model.User, error) {
	return s.userRepo.FindUserByID(id)
}

func (s *userService) CreateUser(requestData request.CreateUser) error {


	// 檢查user 是否已經存在

	personInfo := model.User{
		FirstName: requestData.FirstName,
		LastName: requestData.LastName,
		Email: requestData.Email,
		Phone: requestData.Phone,
	}

	person, _ := s.userRepo.Search(personInfo)


	// 檢查目標銀行是否存在
	bank, err := s.bankRepo.GetBankByID(requestData.BankID)

	if err != nil {
		return err
	}

	if bank == nil {
		return errors.New("bank is not exist")
	}


	tx, err := database.MariaDB.Begin()

	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()


	var userData model.User

	if person == nil {
		userData.FirstName = requestData.FirstName
		userData.LastName = requestData.LastName
		userData.Email = requestData.Email
		userData.Phone = requestData.Phone
		userData.BirthDate = requestData.BirthDate

		userID, err := s.userRepo.InsertUser(tx, userData)
		if err != nil {
			return err
		}


		userData.ID = int(userID)
		 
	} else {
		userData = *person
	}


	account := &model.Account{
		UserID:      userData.ID,
		BankID:      requestData.BankID,
		AccountName: bank.Code + pkg.BuildBankCode(),
		Balance:     100000,
	}

	if err := s.accountRepo.Create(tx, account); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
