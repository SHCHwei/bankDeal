package database


import (
	"fmt"
    "strings"
	"math/rand/v2"
    "bankDeal/internal/model"
	"github.com/brianvoe/gofakeit/v7"
)



func FactoryFake() {

	err := gofakeit.Seed(0)
	if err != nil {
		fmt.Printf("failed to seed gofakeit: %v\n", err)
		return
	}

    var sum int =  1

    var banks []model.Bank

    var bankInsert string = "insert into banks ( code, bankName, capitalAmount) values "

    var bankCodeList string = ""

    for sum <= 10 {
        var bank model.Bank
        var codeSituation bool = true

        err := gofakeit.Struct(&bank)
        if err != nil {
            fmt.Printf("failed to generate fake bank: %v\n", err)
            return
        }

        bank.ID = sum

        for codeSituation {

            code := buildBankCode()
            
            if  codeSituation = strings.Contains(bankCodeList, code) ; !codeSituation {
                bank.Code = code
            }

        }
        
        bankCodeList = bankCodeList + bank.Code + ","

        insertString := fmt.Sprintf("( '%s' , '%s' , %d )," , bank.Code, bank.BankName, bank.CapitalAmount)

        bankInsert = bankInsert + insertString

        banks = append(banks, bank)

        sum += 1
    }

    bankInsert = strings.TrimRight(bankInsert, ",") + ";"

    _, bankErr := MariaDB.Exec(bankInsert)

    if bankErr != nil {
        fmt.Printf("banks is error: %v \n string : %s", bankErr, bankInsert)
    } else {
        fmt.Println("banks is done")
    }


    var persons []model.User
    var userInsert string = "insert into users (FirstName, LastName, Email, Phone) values "
    sum = 1

	for sum <= 100 {
		var person model.User
		err := gofakeit.Struct(&person)
		if err != nil {
			fmt.Printf("failed to generate fake user: %v\n", err)
			return
		}
        person.ID = sum
		person.Phone = "09" + person.Phone

        insertString := fmt.Sprintf("( '%s' , '%s' , '%s' , '%s' ),", person.FirstName, person.LastName, person.Email, person.Phone)

        userInsert = userInsert + insertString

		persons = append(persons, person)
		sum += 1
		
	}

    userInsert = strings.TrimRight(userInsert, ",") + ";"

    _, userErr := MariaDB.Exec(userInsert)

    if userErr != nil {
        fmt.Printf("users is error: %v \n string : %s", userErr, userInsert)
    } else {
        fmt.Println("users is done")
    }



    r := rand.New(rand.NewPCG(1, 2))
    var accountInsert string = "insert into accounts (userID, bankID, accountName , balance) values "

    for _, v := range persons {
        var account model.Account

        err := gofakeit.Struct(&account)
        if err != nil {
            fmt.Printf("failed to generate fake account: %v\n", err)
            return
        }

        targetBank := banks[r.IntN(9)]

        account.UserID = v.ID
        account.BankID = targetBank.ID

        account.AccountName = targetBank.Code + gofakeit.UUID()

        insertString := fmt.Sprintf("( '%d' , '%d' , '%s' , 100000 ),", account.UserID, account.BankID, account.AccountName)

        accountInsert = accountInsert + insertString

    }


    accountInsert = strings.TrimRight(accountInsert, ",") + ";"

    _, accountErr := MariaDB.Exec(accountInsert)

    if accountErr != nil {
        fmt.Printf("accounts is error: %v \n string : %s", accountErr, accountInsert)
    } else {
        fmt.Println("accounts is done")
    }

    
}

func buildBankCode() string{

	// 定義包含大小寫英文字母的字元集
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	length := 3

	// 產生 3 位數的隨機字串
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.IntN(len(letters))]
	}

	return string(b)
}