package model

import (
    "time"
    "database/sql"


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
    CreateBank(tx *sql.Tx, bankData Bank)(int64, error)
}


