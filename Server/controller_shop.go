package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func GetShops() []Shop {
	var shop Shop
	db := ConnectToDatabase()
	defer db.Close()

	sqlStatement := `SELECT * FROM Shop;`
	shops, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer shops.Close()

	var allShops []Shop
	for shops.Next() {
		err = shops.Scan(&shop.ID, &shop.Name, &shop.Location, &shop.AdminUsername, &shop.AdminPassword)
		if err != nil {
			log.Fatal(err)
		}
		allShops = append(allShops, shop)
	}
	return allShops
}

func GetShop(id int) (*Shop, *Error) {
	var shop Shop
	db := ConnectToDatabase()
	defer db.Close()

	sqlStatement := `SELECT * FROM Shop WHERE ID = $1;`
	err := db.QueryRow(sqlStatement, id).Scan(&shop.ID, &shop.Name, &shop.Location, &shop.AdminUsername, &shop.AdminPassword)
	if err != nil {
		return nil, &Error{Status: 404, Error: "This ID doesn't exist"}
	}
	return &shop, nil
}

func CreateShop(r *http.Request) (*Shop, *Error) {
	var shop Shop
	err := json.NewDecoder(r.Body).Decode(&shop)
	if err != nil {
		return nil, &Error{Status: 400, Error: "Invalid Data"}
	}
	db := ConnectToDatabase()
	defer db.Close()

	//var id int
	sqlStatement := `
	INSERT INTO Shop (Name, Location, AdminUsername, AdminPassword)
	VALUES ($1, $2, $3, $4) RETURNING ID, Name, Location, AdminUsername, AdminPassword`
	err = db.QueryRow(sqlStatement, shop.Name, shop.Location, shop.AdminUsername, shop.AdminPassword).Scan(&shop.ID, &shop.Name, &shop.Location, &shop.AdminUsername, &shop.AdminPassword)
	if err != nil {
		log.Fatal(err)
		return nil, &Error{Status: 500, Error: "Error Creating Data"}
	}

	// sqlStatement = `SELECT * FROM "User" WHERE ID = $1;`
	// _ = db.QueryRow(sqlStatement, id).Scan(&user.ID, &user.Name, &user.Age)

	return &shop, nil
}

func UpdateShop(id int, r *http.Request) (*Shop, *Error) {
	db := ConnectToDatabase()
	defer db.Close()

	var (
		shop Shop
		temp Shop
	)
	err := json.NewDecoder(r.Body).Decode(&shop)

	sqlStatement := `SELECT * FROM Shop WHERE ID = $1;`
	err = db.QueryRow(sqlStatement, id).Scan(&temp.ID, &temp.Name, &temp.Location, &temp.AdminUsername, &temp.AdminPassword)
	if err != nil {
		return nil, &Error{Status: 404, Error: "This ID doesn't exist"}
	}

	sqlStatement = `
		UPDATE Shop 
		SET Name = $2, Location = $3, AdminUsername = $4, AdminPassword = $5
		WHERE id = $1;`
	_, err = db.Exec(sqlStatement, id, shop.Name, shop.Location, shop.AdminUsername, shop.AdminPassword)
	if err != nil {
		return nil, &Error{Status: 400, Error: "Invalid Data"}
	}

	shop.ID = id
	return &shop, nil
}

func DeleteShop(id int) (*Shop, *Error) {
	db := ConnectToDatabase()
	defer db.Close()

	var shop Shop
	sqlStatement := `SELECT * FROM Shop WHERE ID = $1;`
	err := db.QueryRow(sqlStatement, id).Scan(&shop.ID, &shop.Name, &shop.Location, &shop.AdminUsername, &shop.AdminPassword)
	if err != nil {
		return nil, &Error{Status: 404, Error: "This ID doesn't exist"}
	}

	sqlStatement = `DELETE FROM Shop WHERE ID = $1;`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatal(err)
		return nil, &Error{Status: 500, Error: "Error Deleting Data"}
	}
	return &shop, nil
}

func ViewPendingOrders(id int) []Order {
	var order Order
	db := ConnectToDatabase()
	defer db.Close()

	sqlStatement := `SELECT * FROM "Order" WHERE ShopID = $1 AND Delivered = false;`
	orders, err := db.Query(sqlStatement, id)
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

func ViewDeliveredOrders(id int) []Order {
	var order Order
	db := ConnectToDatabase()
	defer db.Close()

	sqlStatement := `SELECT * FROM "Order" WHERE ShopID = $1 AND Delivered = true;`
	orders, err := db.Query(sqlStatement, id)
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
