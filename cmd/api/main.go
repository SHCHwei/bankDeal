package main

import (
	"bankDeal/internal/config"
	"bankDeal/internal/database"
	"bankDeal/internal/handler"
	"bankDeal/internal/repository"
	"bankDeal/internal/service"
	"log"
	"net/http"
)

func main() {

	cfg := config.Load()
	database.BuildTables()

	if dbResult := database.MariaDBConnect(*cfg); dbResult == false {
		log.Fatal("ＳＯＲＲＹ!! Database is Connect Failed ")
	}

	// build Repository
	userRepo := repository.NewUserRepository()
	accountRepo := repository.NewAccountRepository()
	dealRepo := repository.NewDealRepositories()
	bankRepo := repository.NewBankRepository()

	// build Service
	userSvc := service.NewUserService(userRepo, accountRepo, bankRepo)
	dealSvc := service.NewDealService(dealRepo, accountRepo)

	handleDeal := handler.NewDealHandler(dealSvc)
	handleUser := handler.NewUserHandler(userSvc)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /deals/{id}", handleDeal.GetDeal) // 查詢單筆交易
	mux.HandleFunc("POST /deals", handleDeal.CreateDeal)  // 建立交易

	mux.HandleFunc("GET /user/{id}", handleUser.GetUser) // 查詢用戶
	mux.HandleFunc("POST /user/create", handleUser.CreateUser) // 建立用戶

	mux.HandleFunc("GET /init/fake", func(w http.ResponseWriter, r *http.Request) {
		database.FactoryFake()
	}) //建立假資料

	cfg.Init = false

	log.Printf("server listening on %s \n", cfg.Addr)

	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatal(err)
	}

}
