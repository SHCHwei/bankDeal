
# 主架構

參考 [golang-standards ](https://github.com/golang-standards/project-layout/blob/master/README_zh-TW.md)

project/  
├── cmd/  
│   └── api/  
│       └── main.go  
├── internal/  
│   ├── handler/        # HTTP handler（接 request）  
│   ├── service/        # 商業邏輯  
│   ├── repository/     # DB 存取  
│   ├── model/          # struct + interface 定義  
│   ├── configs/        # 設定檔  
│   ├── database/       # 資料庫連接  
│   │   ├── mariadb.go  # MariaDB連接設定    
│   │   └── seed.go     # 假資料輸入功能  
│   └── middleware/     # 中介層  
├── pkg/                # 可重用工具（對外）  
├── migrations/         # SQL檔  
├── go.mod  
├── go.sum  
├── .ari.toml           # Go air-verse 套件設定檔  
├── docker-compose.yml  
└── Dockerfile  




---

# 主要依賴套件

套件列表:  
    * golang-migrate/migrate/v4  
    * air-verse/air@v1.52  


---