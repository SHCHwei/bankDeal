package main

import (
	"bankDeal/internal/config"
	"bankDeal/internal/database"
	"bankDeal/internal/handler"
	"bankDeal/internal/logger"
	"bankDeal/internal/repository"
	"bankDeal/internal/service"
	"bankDeal/internal/swagger"
	"bankDeal/internal/middleware"
	"log"
	"net/http"
)

func main() {

	cfg, dbCfg := config.Load()
	database.BuildTables()

	if dbResult := database.MariaDBConnect(*dbCfg); !dbResult {
		log.Fatal("ＳＯＲＲＹ!! Database is Connect Failed ")
	}

	// Initialize logger
	log_, err := logger.GetInstance()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer log_.Close()

	// build Repository
	userRepo := repository.NewUserRepository(database.MariaDB)
	accountRepo := repository.NewAccountRepository(database.MariaDB)
	dealRepo := repository.NewDealRepositories(database.MariaDB)
	bankRepo := repository.NewBankRepository(database.MariaDB)

	// build Service
	userSvc := service.NewUserService(database.MariaDB, userRepo, accountRepo, bankRepo)
	dealSvc := service.NewDealService(database.MariaDB, dealRepo, accountRepo)

	handleDeal := handler.NewDealHandler(dealSvc)
	handleUser := handler.NewUserHandler(userSvc)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /deals/{id}", handleDeal.GetDeal)   // 查詢單筆交易
	mux.HandleFunc("POST /deals", handleDeal.CreateDeal) // 建立交易

	mux.HandleFunc("GET /user/{id}", handleUser.GetUser)          // 查詢用戶
	mux.HandleFunc("POST /user/create", handleUser.CreateUser) // 建立用戶

	mux.HandleFunc("GET /init/fake", func(w http.ResponseWriter, r *http.Request) {
		database.FactoryFake()
		w.WriteHeader(http.StatusNoContent)
	}) //建立假資料

	mux.HandleFunc("GET /swagger", swagger.UIHandler)
	mux.HandleFunc("GET /swagger/doc", swagger.DocHandler)

	handler := middleware.Chain(
		mux,
		middleware.Logging,
		// middleware.Auth,
	)

	log.Printf("server listening on %s \n", cfg.Addr)

	if err := http.ListenAndServe(cfg.Addr, handler); err != nil {
		log.Fatal(err)
	}

}
