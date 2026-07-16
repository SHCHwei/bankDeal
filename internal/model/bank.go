package model

import (
    "time"
    "context"
)

type Bank struct {
    ID                  int
    Code                string
    BankName            string      `fake:"{appname}"`
    CapitalAmount       uint32      `fake:"{number:100000, 99999999}"`
   	CreatedAt           time.Time   `fake:"-"`
   	UpdatedAt           time.Time   `fake:"-"`
}



type BankRepository interface{
    GetBankByID(id int)(*Bank, error)
    CreateBank(ctx context.Context, bankData Bank)(int64, error)
}


