package request

type CreateUser struct {
		FirstName string    `schema:"first_name"`
		LastName  string    `schema:"last_name"`
		Email     string    `schema:"email"`
		Phone     string    `schema:"phone"`
		BirthDate string    `schema:"birthdate"`
		BankID	  int		`schema:"bank_ID"`
}