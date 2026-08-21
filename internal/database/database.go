package database


// AI 建立

import (
	"context"
	"database/sql"
)


type contextKey struct{}
var txKey = contextKey{}

type TxManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

// Transaction 在此方法中執行，傳入一個閉包 (closure)
func (m *TxManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// 將 tx 塞入 context 中，傳遞給下游
	txCtx := context.WithValue(ctx, txKey, tx)

	// 執行業務邏輯
	err = fn(txCtx)
	if err != nil {
		roll_err := tx.Rollback() // 出錯就 Rollback
		if roll_err != nil {
			return roll_err
		}
		return err
	}

	return tx.Commit() // 成功就 Commit
}

// 提供一個輔助函式，讓 Repository 取得當前的執行器 (sql.Tx 或 sql.DB)
func GetExecutor(ctx context.Context, defaultDB *sql.DB) Executor {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx
	}
	return defaultDB
}

// 定義一個介面，因為 *sql.DB 與 *sql.Tx 都實作了這些方法
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

