package main

import "log"

type Customer struct {
	ID                   int    `json:"id"`
	Email                string `json:"email"`
	FirstName            string `json:"firstName"`
	LastName             string `json:"lastName"`
	Password             string `json:"password"`
	CreditCardNumber     string `json:"creditCardNumber"`
	CreditCardExpiryDate string `json:"creditCardExpiryDate"`
	CreditCardCVV        int    `json:"creditCardCVV"`
}

func CreateCustomerModel() {
	db := ConnectToDatabase()
	defer db.Close()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Customer (
		ID SERIAL PRIMARY KEY,
		Email VARCHAR(255) UNIQUE NOT NULL,
		FirstName VARCHAR(255) NOT NULL,
		LastName VARCHAR(255) NOT NULL,
		Password VARCHAR(255) NOT NULL,
		CreditCardNumber VARCHAR(255) UNIQUE NOT NULL,
		CreditCardExpiryDate DATE NOT NULL,
		CreditCardCVV INT NOT NULL
		);
	`)

	if err != nil {
		log.Fatal(err)
	}
}
