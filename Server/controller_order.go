package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func GetOrders() []Order {
	var order Order
	db := ConnectToDatabase()
	defer db.Close()

	sqlStatement := `SELECT * FROM "Order";`
	orders, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer orders.Close()

	var allOrders []Order
	for orders.Next() {
		err = orders.Scan(&order.ID, &order.Delivered, &order.Price, &order.Date, &order.CustomerID, &order.ShopID)
		if err != nil {
			log.Fatal(err)
		}
		allOrders = append(allOrders, order)
	}
	return allOrders
}

func GetOrder(id int) (*Order, *Error) {
	var order Order
	db := ConnectToDatabase()
	defer db.Close()

	sqlStatement := `SELECT * FROM "Order" WHERE ID = $1;`
	err := db.QueryRow(sqlStatement, id).Scan(&order.ID, &order.Delivered, &order.Price, &order.Date, &order.CustomerID, &order.ShopID)
	if err != nil {
		return nil, &Error{Status: 404, Error: "This ID doesn't exist"}
	}
	return &order, nil
}

func CreateOrder(r *http.Request) (*Order, *Error) {
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		return nil, &Error{Status: 400, Error: "Invalid Data"}
	}
	db := ConnectToDatabase()
	defer db.Close()

	var sqlStatement string
	if order.Date == "" {
		sqlStatement = `
		INSERT INTO "Order" (Delivered, Price, CustomerID, ShopID)
		VALUES ($1, $2, $3, $4) RETURNING ID, Delivered, Price, Date, CustomerID, ShopID`
		err = db.QueryRow(sqlStatement, order.Delivered, order.Price, order.CustomerID, order.ShopID).Scan(&order.ID, &order.Delivered, &order.Price, &order.Date, &order.CustomerID, &order.ShopID)
	} else {
		sqlStatement = `
		INSERT INTO "Order" (Delivered, Price, Date, CustomerID, ShopID)
		VALUES ($1, $2, $3, $4, $5) RETURNING ID, Delivered, Price, Date, CustomerID, ShopID`
		err = db.QueryRow(sqlStatement, order.Delivered, order.Price, order.Date, order.CustomerID, order.ShopID).Scan(&order.ID, &order.Delivered, &order.Price, &order.Date, &order.CustomerID, &order.ShopID)
	}
	if err != nil {
		log.Fatal(err)
		return nil, &Error{Status: 500, Error: "Error Creating Data"}
	}

	return &order, nil
}

func UpdateOrder(id int, r *http.Request) (*Order, *Error) {
	db := ConnectToDatabase()
	defer db.Close()

	var (
		order Order
		temp  Order
	)
	err := json.NewDecoder(r.Body).Decode(&order)

	sqlStatement := `SELECT * FROM "Order" WHERE ID = $1;`
	err = db.QueryRow(sqlStatement, id).Scan(&temp.ID, &temp.Delivered, &temp.Price, &temp.Date, &temp.CustomerID, &temp.ShopID)
	if err != nil {
		return nil, &Error{Status: 404, Error: "This ID doesn't exist"}
	}

	sqlStatement = `
		UPDATE "Order" 
		SET Delivered = $2, Price = $3, Date = $4, CustomerID = $5, ShopID = $6
		WHERE id = $1;`
	_, err = db.Exec(sqlStatement, id, order.Delivered, order.Price, order.Date, order.CustomerID, order.ShopID)
	if err != nil {
		return nil, &Error{Status: 400, Error: "Invalid Data"}
	}
	order.ID = id
	return &order, nil
}

func DeleteOrder(id int) (*Order, *Error) {
	db := ConnectToDatabase()
	defer db.Close()
	var order Order
	sqlStatement := `SELECT * FROM "Order" WHERE ID = $1;`
	err := db.QueryRow(sqlStatement, id).Scan(&order.ID, &order.Delivered, &order.Price, &order.Date, &order.CustomerID, &order.ShopID)
	if err != nil {
		return nil, &Error{Status: 404, Error: "This ID doesn't exist"}
	}

	sqlStatement = `DELETE FROM "Order" WHERE ID = $1;`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatal(err)
		return nil, &Error{Status: 500, Error: "Error Deleting Data"}
	}
	return &order, nil
}

func ViewOrderItems(orderID int) ([]Item, *Error) {
	var item Item
	db := ConnectToDatabase()
	defer db.Close()

	_, orderReadError := GetOrder(orderID)
	if orderReadError != nil {
		return nil, orderReadError
	}

	sqlStatement := `
		SELECT I.ID, I.Name, I.Type, I.Price, I.Description, I.ImageURL, I.ShopID
		FROM Item I INNER JOIN Contain C ON I.ID = C.ItemID
		WHERE C.OrderID = $1;
	`
	items, err := db.Query(sqlStatement, orderID)
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
	return allItems, nil

}
