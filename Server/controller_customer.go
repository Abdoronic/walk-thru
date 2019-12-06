package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/token"
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

func Checkout(customerID int, orderID int, shopID int, r *http.Request) (*Order, *Error) {
	stripe.Key = GetConfig().StripeKey

	order, orderError := GetOrder(orderID)
	if orderError != nil {
		return nil, orderError
	}

	_, shopError := GetShop(shopID)
	if shopError != nil {
		return nil, shopError
	}

	customer, customerError := GetCustomer(customerID)
	if customerError != nil {
		return nil, customerError
	}

	if int(order.CustomerID) != customerID {
		return nil, &Error{Status: 400, Error: "Order doesn't belong to customer"}
	}

	db := ConnectToDatabase()
	defer db.Close()

	var (
		orderItemIDs []int
		shopItemIDs  []int
	)

	sqlStatement := `
		SELECT ItemID FROM Contain
		WHERE OrderID = $1
	`
	orderItemIDsIterator, orderReadError := db.Query(sqlStatement, orderID)
	if orderReadError != nil {
		log.Fatal(orderReadError)
	}
	defer orderItemIDsIterator.Close()

	for orderItemIDsIterator.Next() {
		var itemID int
		err := orderItemIDsIterator.Scan(&itemID)
		if err != nil {
			log.Fatal(err)
		}
		orderItemIDs = append(orderItemIDs, itemID)
	}

	sqlStatement = `
		SELECT ID FROM Item
		WHERE ShopID = $1
	`
	shopItemIDsIterator, shopReadError := db.Query(sqlStatement, shopID)
	if shopReadError != nil {
		log.Fatal(shopReadError)
	}
	defer shopItemIDsIterator.Close()

	for shopItemIDsIterator.Next() {
		var itemID int
		err := shopItemIDsIterator.Scan(&itemID)
		if err != nil {
			log.Fatal(err)
		}
		shopItemIDs = append(shopItemIDs, itemID)
	}

	for i := 0; i < len(orderItemIDs); i++ {
		contain := false
		for j := 0; j < len(shopItemIDs); j++ {
			if orderItemIDs[j] == shopItemIDs[i] {
				contain = true
				break
			}
		}
		if !contain {
			return nil, &Error{Status: 400, Error: "Order Items must be from same shop"}
		}
	}

	//Start Payment
	cardParams := &stripe.TokenParams{
		Card: &stripe.CardParams{
			Number:   stripe.String(strconv.Itoa(customer.CreditCardNumber)),
			ExpMonth: stripe.String(customer.CreditCardExpiryDate[:2]),
			ExpYear:  stripe.String(customer.CreditCardExpiryDate[2:]),
			CVC:      stripe.String(strconv.Itoa(customer.CreditCardCVV))},
		// Number:   stripe.String("4242424242424242"),
		// ExpMonth: stripe.String("12"),
		// ExpYear:  stripe.String("20"),
		// CVC:      stripe.String("123")},
	}
	cardToken, err := token.New(cardParams)

	if err != nil {
		fmt.Println(err)
		return nil, &Error{Status: 500, Error: "Invaild credit card info"}
	}

	params := &stripe.ChargeParams{
		Amount:      stripe.Int64(int64(order.Price * 100)),
		Currency:    stripe.String(string(stripe.CurrencyEGP)),
		Description: stripe.String("Order Payment"),
	}
	params.SetSource(cardToken.ID)

	_, err = charge.New(params)

	if err != nil {
		fmt.Println(err)
		return nil, &Error{Status: 500, Error: "Failed Transaction"}
	}

	//Mark order as recieved
	sqlStatement = `
		UPDATE "Order" 
		SET ShopID = $2
		WHERE id = $1;`
	_, err = db.Exec(sqlStatement, orderID, shopID)
	if err != nil {
		return nil, &Error{Status: 400, Error: "Couldn't mark order as recieved"}
	}

	order, _ = GetOrder(orderID)
	return order, nil
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

// As a Customer i can view my Order Items.
type OrderItem struct {
	OrderID     int           `json:"orderID"`
	ItemID      int           `json:"itemID"`
	CustomerID  int           `json:"customerID"`
	ShopID      sql.NullInt64 `json:"shopID"`
	Quantity    int           `json:"quantity"`
	Delivered   bool          `json:"delivered"`
	OrderPrice  float64       `json:"OrderPrice"`
	Date        string        `json:"date"`
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Price       float64       `json:"price"`
	Description string        `json:"description"`
	ImageURL    string        `json:"imageURL"`
	ItemShop    int           `json:"itemShop"`
}

func CustomerViewOrderItems(customerID int, orderID int, r *http.Request) ([]OrderItem, *Error) {
	var orderItem OrderItem
	db := ConnectToDatabase()
	defer db.Close()

	getCustomerSQL := `SELECT * FROM Customer WHERE ID = $1;`
	var customer Customer
	customerErr := db.QueryRow(getCustomerSQL, customerID).Scan(&customer.ID, &customer.Email, &customer.FirstName, &customer.LastName, &customer.CreditCardNumber, &customer.CreditCardExpiryDate, &customer.CreditCardCVV)

	sqlStatement := `SELECT * FROM "Order" WHERE ID = $1;`
	var order Order
	orderErr := db.QueryRow(sqlStatement, orderID).Scan(&order.ID, &order.Delivered, &order.Price, &order.Date, &order.CustomerID, &order.ShopID)

	if customerErr != nil || orderErr != nil {
		return nil, &Error{Status: 404, Error: "This ID doesn't exist"}
	}

	viewOrderItemsSQL := `SELECT Contain.OrderID, Contain.ItemID, "Order".CustomerID, "Order".ShopID, Contain.Quantity, "Order".Delivered, "Order".Price AS OrderPrice, "Order".Date, Item.Name, Item.Type, Item.Price, Item.Description, Item.ImageURL, Item.ShopID AS ItemShop
	FROM Contain, "Order", Item
	WHERE "Order".ID = Contain.OrderID AND Item.ID = Contain.ItemID AND Contain.OrderID = $1 AND "Order".CustomerID = $2
	`
	noResult := db.QueryRow(viewOrderItemsSQL, orderID, customerID).Scan()
	if noResult == sql.ErrNoRows {
		return nil, &Error{Status: 400, Error: "This order is Empty or it does't belong to you"}
	}
	orderItems, viewOrderItemsErr := db.Query(viewOrderItemsSQL, orderID, customerID)
	if viewOrderItemsErr != nil {
		log.Fatal(viewOrderItemsErr)
	}
	defer orderItems.Close()
	var allOrderItems []OrderItem
	for orderItems.Next() {
		err := orderItems.Scan(&orderItem.OrderID, &orderItem.ItemID, &orderItem.CustomerID, &orderItem.ShopID, &orderItem.Quantity, &orderItem.Delivered, &orderItem.OrderPrice, &orderItem.Date, &orderItem.Name, &orderItem.Type, &orderItem.Price, &orderItem.Description, &orderItem.ImageURL, &orderItem.ItemShop)
		if err != nil {
			log.Fatal(err)
		}
		allOrderItems = append(allOrderItems, orderItem)
	}
	return allOrderItems, nil
}
