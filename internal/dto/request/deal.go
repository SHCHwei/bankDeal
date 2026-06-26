package request

type CreateDeal struct {
	AccountID        int    `schema:"account_id" validate:"required"`
	Volume           int64  `schema:"volume" validate:"required,min=1,max=50000"`
	TransactionType  uint8  `schema:"transaction_type"`
	TradingAccountID int    `schema:"trading_account_id"`
	Remark           string `schema:"remark"`
}
