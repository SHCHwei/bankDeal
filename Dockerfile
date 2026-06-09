FROM golang:1.23-alpine

# 安裝 git (air 編譯時可能需要)
RUN apk add --no-cache git

# 設定工作目錄
WORKDIR /app

# 先複製 go mod（加速 build cache）
COPY go.* ./
#RUN go mod download

# 複製全部程式碼
COPY . .

# 安裝 air
RUN go install github.com/air-verse/air@v1.52

# 開放 port（依你的服務調整）
EXPOSE 8080

# 啟動 air
CMD ["air"]