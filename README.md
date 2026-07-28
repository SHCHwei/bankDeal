# bankDeal

`bankDeal` 是一個以 Go 開發的銀行交易管理 API，支援使用者、帳戶與交易紀錄的建立與查詢。  
目標是不使用框架，只利用原生的語法建制，期望過程中可以加深了解GO語言的特性。  
並且搭配 bankDeal-client-test 使用  
https://github.com/SHCHwei/bankDeal-client-test

## 主要功能

- 使用者建立與查詢
- 交易建立與查詢
- Swagger API 文件
- 資料庫初始化與假資料產生
- Docker 與 docker-compose 支援

## 專案結構

- `cmd/api/`：應用程式啟動入口
- `internal/config/`：設定載入與環境讀取
- `internal/database/`：資料庫連線、建立表格與測試資料
- `internal/dto/`：HTTP 請求資料結構
- `internal/handler/`：HTTP 路由處理程式
- `internal/repository/`：資料存取層
- `internal/service/`：業務邏輯層
- `internal/middleware/`：日誌、CORS、Auth 等中介層
- `internal/swagger/`：Swagger UI 與文件路由
- `migrations/`：資料庫 schema SQL 檔案

## 執行環境

- Go 1.23
- MySQL / MariaDB
- Docker (選用)

## 啟動方式

### 1. 直接本機執行

1. 複製專案
2. 檢查 `.env.example` 內容，必要時調整設定
3. 執行資料庫
4. 啟動應用程式：

```bash
go run ./cmd/api
```

預設會讀取 `.env.example`，並使用以下預設值：

- `addr`: `:8080`
- `user`: `user`
- `password`: `user1`
- `database_host`: `bankdeal-db-1`
- `port`: `3306`
- `database`: `bankDeal`

### 2. 使用 Docker Compose

```bash
docker-compose up --build
```

`docker-compose.yml` 會啟動兩個服務：

- `app`：Go 應用程式，對外開放 `8080`
- `db`：MariaDB

## API 路由

- `GET /deals/{id}`：查詢單筆交易
- `POST /deals`：建立交易
- `GET /user/{id}`：查詢使用者
- `POST /user/create`：建立使用者
- `GET /init/fake`：產生假資料
- `GET /swagger`：Swagger UI
- `GET /swagger/doc`：Swagger JSON 文件

## 資料庫遷移

專案目錄中的 `migrations/` 包含初始資料表建立 SQL：

- `001_init_banks.up.sql`
- `002_init_users.up.sql`
- `003_init_accounts.up.sql`
- `004_init_deals.up.sql`

## 測試

執行單元測試：

```bash
go test ./...
```

## 開發建議

- 如果需要額外的設定，建議新增 `.env` 或修改 `.env.example`
- 若要啟用授權，可在 `cmd/api/main.go` 中解除 `middleware.Auth` 註解
- Swagger UI 可協助檢視 API 結構與請求範例

## 貢獻

歡迎提交 Issue 或 Pull Request，並依照既有專案風格新增功能或修正問題。
