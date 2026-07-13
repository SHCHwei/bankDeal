package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Logger 日誌記錄器
type Logger struct {
	file *os.File
	mu   sync.Mutex
}

var (
	instance *Logger
	once     sync.Once
)

// GetInstance 取得日誌記錄器單例
func GetInstance() (*Logger, error) {
	var err error
	once.Do(func() {
		instance, err = NewLogger()
	})
	return instance, err
}

// NewLogger 建立新的日誌記錄器
func NewLogger() (*Logger, error) {
	// 建立 logs 目錄（如果不存在）
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("建立日誌目錄失敗: %w", err)
	}

	// 建立日誌檔案
	logFile := filepath.Join(logsDir, "deals.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("開啟日誌檔案失敗: %w", err)
	}

	return &Logger{file: file}, nil
}

// LogCreateDeal 記錄建立交易的日誌
func (l *Logger) LogCreateDeal(
	accountID int,
	volume int64,
	transactionType uint8,
	tradingAccountID int,
	remark string,
	success bool,
	dealID int,
	errorMsg string,
) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	status := "成功"
	if !success {
		status = "失敗"
	}

	logEntry := fmt.Sprintf(
		"[%s] 操作: 建立交易 | 狀態: %s | 帳戶ID: %d | 交易金額: %d | 交易類型: %d | 對方帳戶: %d | 備註: %s | 交易ID: %d",
		timestamp,
		status,
		accountID,
		volume,
		transactionType,
		tradingAccountID,
		remark,
		dealID,
	)

	if !success && errorMsg != "" {
		logEntry += fmt.Sprintf(" | 錯誤: %s", errorMsg)
	}

	logEntry += "\n"

	if _, err := l.file.WriteString(logEntry); err != nil {
		fmt.Fprintf(os.Stderr, "寫入日誌檔案失敗: %v\n", err)
	}
}

// LogRequestError 記錄請求錯誤
func (l *Logger) LogRequestError(
	accountID int,
	volume int64,
	transactionType uint8,
	tradingAccountID int,
	errorType string,
	errorMsg string,
) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	logEntry := fmt.Sprintf(
		"[%s] 操作: 建立交易 | 狀態: 請求錯誤 | 帳戶ID: %d | 交易金額: %d | 交易類型: %d | 對方帳戶: %d | 錯誤類型: %s | 錯誤訊息: %s\n",
		timestamp,
		accountID,
		volume,
		transactionType,
		tradingAccountID,
		errorType,
		errorMsg,
	)

	if _, err := l.file.WriteString(logEntry); err != nil {
		fmt.Fprintf(os.Stderr, "寫入日誌檔案失敗: %v\n", err)
	}
}

// Close 關閉日誌文件
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// GetLogFilePath 取得日誌檔案路徑
func GetLogFilePath() string {
	return filepath.Join("logs", "deals.log")
}
