package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func GetCustomers() []Customer {
	var customer Customer
	db := ConnectToDatabase()
	defer db.Close()

	sqlStatement := `SELECT * FROM Customer;`
	customers, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer customers.Close()

	var allCustomers []Customer
	for customers.Next() {
		err = customers.Scan(&customer.ID, &customer.Email, &customer.FirstName, &customer.LastName, &customer.CreditCardNumber, &customer.CreditCardExpiryDate, &customer.CreditCardCVV)
		if err != nil {
			log.Fatal(err)
		}
		allCustomers = append(allCustomers, customer)
	}
	return allCustomers
}

func GetCustomer(id int) (*Customer, *Error) {
	var customer Customer
	db := ConnectToDatabase()
	defer db.Close()

	sqlStatement := `SELECT * FROM Customer WHERE ID = $1;`
	err := db.QueryRow(sqlStatement, id).Scan(&customer.ID, &customer.Email, &customer.FirstName, &customer.LastName, &customer.CreditCardNumber, &customer.CreditCardExpiryDate, &customer.CreditCardCVV)
	if err != nil {
		return nil, &Error{Status: 404, Error: "This ID doesn't exist"}
	}
	return &customer, nil
}

func CreateCustomer(r *http.Request) (*Customer, *Error) {
	var customer Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		return nil, &Error{Status: 400, Error: "Invalid Data"}
	}
	db := ConnectToDatabase()
	defer db.Close()

	//var id int
	sqlStatement := `
	INSERT INTO Customer (Email, FirstName, LastName, CreditCardNumber, CreditCardExpiryDate, CreditCardCVV)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING ID, Email, FirstName, LastName, CreditCardNumber, CreditCardExpiryDate, CreditCardCVV`
	err = db.QueryRow(sqlStatement, customer.Email, customer.FirstName, customer.LastName, customer.CreditCardNumber, customer.CreditCardExpiryDate, customer.CreditCardCVV).Scan(&customer.ID, &customer.Email, &customer.FirstName, &customer.LastName, &customer.CreditCardNumber, &customer.CreditCardExpiryDate, &customer.CreditCardCVV)
	if err != nil {
		log.Fatal(err)
		return nil, &Error{Status: 500, Error: "Error Creating Data"}
	}

	// sqlStatement = `SELECT * FROM "User" WHERE ID = $1;`
	// _ = db.QueryRow(sqlStatement, id).Scan(&user.ID, &user.Name, &user.Age)

	return &customer, nil
}

func UpdateCustomer(id int, r *http.Request) (*Customer, *Error) {
	db := ConnectToDatabase()
	defer db.Close()

	var (
		customer Customer
		temp     Customer
	)
	err := json.NewDecoder(r.Body).Decode(&customer)

	sqlStatement := `SELECT * FROM Customer WHERE ID = $1;`
	err = db.QueryRow(sqlStatement, id).Scan(&temp.ID, &temp.Email, &temp.FirstName, &temp.LastName, &temp.CreditCardNumber, &temp.CreditCardExpiryDate, &temp.CreditCardCVV)
	if err != nil {
		return nil, &Error{Status: 404, Error: "This ID doesn't exist"}
	}

	sqlStatement = `
		UPDATE Customer 
		SET Email = $2, FirstName = $3, LastName = $4, CreditCardNumber = $5, CreditCardExpiryDate = $6, CreditCardCVV = $7
		WHERE id = $1;`
	_, err = db.Exec(sqlStatement, id, customer.Email, customer.FirstName, customer.LastName, customer.CreditCardNumber, customer.CreditCardExpiryDate, customer.CreditCardCVV)
	if err != nil {
		return nil, &Error{Status: 400, Error: "Invalid Data"}
	}

	// sqlStatement = `SELECT * FROM "User" WHERE ID = $1;`
	// _ = db.QueryRow(sqlStatement, id).Scan(&user.ID, &user.Name, &user.Age)
	customer.ID = id
	return &customer, nil
}

func DeleteCustomer(id int) (*Customer, *Error) {
	db := ConnectToDatabase()
	defer db.Close()

	var customer Customer
	sqlStatement := `SELECT * FROM Customer WHERE ID = $1;`
	err := db.QueryRow(sqlStatement, id).Scan(&customer.ID, &customer.Email, &customer.FirstName, &customer.LastName, &customer.CreditCardNumber, &customer.CreditCardExpiryDate, &customer.CreditCardCVV)
	if err != nil {
		return nil, &Error{Status: 404, Error: "This ID doesn't exist"}
	}

	sqlStatement = `DELETE FROM Customer WHERE ID = $1;`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatal(err)
		return nil, &Error{Status: 500, Error: "Error Deleting Data"}
	}
	return &customer, nil
}

func ViewItems(id int) []Item {
	var item Item
	db := ConnectToDatabase()
	defer db.Close()

	sqlStatement := `SELECT * FROM Item WHERE ShopID = $1;`
	items, err := db.Query(sqlStatement, id)
	if err != nil {
		log.Fatal(err)
	}
	defer items.Close()

	var allItems []Item
	for items.Next() {
		err = items.Scan(&item.ID, &item.Name, &item.Type, &item.Price, &item.Description, &item.ImageURL, &item.ShopID)
		if err != nil {
			log.Fatal(err)
		}
		allItems = append(allItems, item)
	}
	return allItems
}
