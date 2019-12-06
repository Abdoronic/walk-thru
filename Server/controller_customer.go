package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
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

	sqlStatement := `
	INSERT INTO Customer (Email, FirstName, LastName, CreditCardNumber, CreditCardExpiryDate, CreditCardCVV)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING ID, Email, FirstName, LastName, CreditCardNumber, CreditCardExpiryDate, CreditCardCVV`
	err = db.QueryRow(sqlStatement, customer.Email, customer.FirstName, customer.LastName, customer.CreditCardNumber, customer.CreditCardExpiryDate, customer.CreditCardCVV).Scan(&customer.ID, &customer.Email, &customer.FirstName, &customer.LastName, &customer.CreditCardNumber, &customer.CreditCardExpiryDate, &customer.CreditCardCVV)
	if err != nil {
		log.Fatal(err)
		return nil, &Error{Status: 500, Error: "Error Creating Data"}
	}

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

// As a Customer i can create an Order.
func CustomerCreateOrder(id int, r *http.Request) (*Order, *Error) {
	db := ConnectToDatabase()
	defer db.Close()

	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		log.Fatal(err)
	}

	order.CustomerID = id
	modifiedBody, err := json.Marshal(order)
	if err != nil {
		log.Fatal(err)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(modifiedBody))
	r.ContentLength = int64(len(modifiedBody))

	createdOrder, createError := CreateOrder(r)
	if createError != nil {
		return nil, createError
	}
	return createdOrder, nil
}

// As a Customer i can add items to my Order.
func CustomerAddItem(orderID int, itemID int, r *http.Request) (*Order, *Error) {
	db := ConnectToDatabase()
	defer db.Close()

	var currentQuantity int
	sqlStatement := `SELECT Quantity FROM Contain WHERE OrderID = $1 AND ItemID = $2;`
	errQuery := db.QueryRow(sqlStatement, orderID, itemID).Scan(&currentQuantity)
	switch errQuery {
	case sql.ErrNoRows:
		{
			// Insert into contain relation
			createContainSQL := `INSERT INTO Contain (OrderID, ItemID, Quantity)
				VALUES ($1, $2, $3)`
			_, errContain := db.Exec(createContainSQL, orderID, itemID, 1)
			if errContain != nil {
				log.Fatal(errContain)
			}
		}
	case nil:
		{
			// update contain relation
			updateContainSQL := `UPDATE Contain
			SET Quantity = $3
			WHERE ItemID = $1 AND OrderID = $2;`
			_, errContain := db.Exec(updateContainSQL, itemID, orderID, currentQuantity+1)
			if errContain != nil {
				log.Fatal(errContain)
			}
		}
	default:
		log.Fatal(errQuery)
	}
	// get item price
	item, errItem := GetItem(itemID)
	if errItem != nil {
		return nil, errItem
	}
	// update order table
	order, errOrder := GetOrder(orderID)
	if errOrder != nil {
		return nil, errOrder
	}
	order.Price = order.Price + item.Price
	modifiedBody, err := json.Marshal(order)
	if err != nil {
		log.Fatal(err)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(modifiedBody))
	r.ContentLength = int64(len(modifiedBody))
	updatedOrder, errOrderUpdate := UpdateOrder(orderID, r)
	if errOrderUpdate != nil {
		return nil, errOrderUpdate
	}
	return updatedOrder, nil
}

// As a Customer i can remove items from my Order.
func CustomerRemoveItem(orderID int, itemID int, r *http.Request) (*Order, *Error) {
	db := ConnectToDatabase()
	defer db.Close()

	var currentQuantity int
	sqlStatement := `SELECT Quantity FROM Contain WHERE OrderID = $1 AND ItemID = $2;`
	errQuery := db.QueryRow(sqlStatement, orderID, itemID).Scan(&currentQuantity)
	if errQuery != nil {
		log.Fatal(errQuery)
	}
	// get item price
	item, errItem := GetItem(itemID)
	if errItem != nil {
		return nil, errItem
	}
	if currentQuantity == 1 {
		// delete from contain relation
		deleteContainSQL := `DELETE FROM Contain
		WHERE ItemID = $1 AND OrderID = $2;`
		_, errContain := db.Exec(deleteContainSQL, itemID, orderID)
		if errContain != nil {
			log.Fatal(errContain)
		}
	} else {
		// update contain relation
		updateContainSQL := `UPDATE Contain
		SET Quantity = $3
		WHERE ItemID = $1 AND OrderID = $2;`
		_, errContain := db.Exec(updateContainSQL, itemID, orderID, currentQuantity-1)
		if errContain != nil {
			log.Fatal(errContain)
		}
	}

	// update order table
	order, errOrder := GetOrder(orderID)
	if errOrder != nil {
		return nil, errOrder
	}
	order.Price = order.Price - item.Price
	modifiedBody, err := json.Marshal(order)
	if err != nil {
		log.Fatal(err)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(modifiedBody))
	r.ContentLength = int64(len(modifiedBody))
	updatedOrder, errOrderUpdate := UpdateOrder(orderID, r)
	if errOrderUpdate != nil {
		return nil, errOrderUpdate
	}
	return updatedOrder, nil
}
