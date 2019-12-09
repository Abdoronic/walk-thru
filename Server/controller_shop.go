package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
)

func GetShops() ([]Shop, *Error) {
	var shop Shop
	db := ConnectToDatabase()
	defer db.Close()

	sqlStatement := `SELECT * FROM Shop;`
	shops, err := db.Query(sqlStatement)
	if err != nil {
		glog.Error(err)
		return nil, nil
	}
	defer shops.Close()

	var allShops []Shop
	for shops.Next() {
		err = shops.Scan(&shop.ID, &shop.Name, &shop.Location, &shop.AdminUsername, &shop.AdminPassword)
		if err != nil {
			glog.Error(err)
			return nil, nil
		}
		allShops = append(allShops, shop)
	}
	if allShops == nil {
		return nil, &Error{Status: 404, Error: "No Shops exist"}
	}
	return allShops, nil
}

func GetShop(id int) (*Shop, *Error) {
	var shop Shop
	db := ConnectToDatabase()
	defer db.Close()

	sqlStatement := `SELECT * FROM Shop WHERE ID = $1;`
	err := db.QueryRow(sqlStatement, id).Scan(&shop.ID, &shop.Name, &shop.Location, &shop.AdminUsername, &shop.AdminPassword)
	if err != nil {
		glog.Error(err)
		return nil, &Error{Status: 404, Error: "This ID doesn't exist"}
	}
	return &shop, nil
}

func CreateShop(r *http.Request) (*Shop, *Error) {
	var shop Shop
	err := json.NewDecoder(r.Body).Decode(&shop)
	if err != nil {
		glog.Error(err)
		return nil, &Error{Status: 400, Error: "Invalid Data"}
	}
	db := ConnectToDatabase()
	defer db.Close()

	sqlStatement := `
	INSERT INTO Shop (Name, Location, AdminUsername, AdminPassword)
	VALUES ($1, $2, $3, $4) RETURNING ID, Name, Location, AdminUsername, AdminPassword`
	err = db.QueryRow(sqlStatement, shop.Name, shop.Location, shop.AdminUsername, shop.AdminPassword).Scan(&shop.ID, &shop.Name, &shop.Location, &shop.AdminUsername, &shop.AdminPassword)
	if err != nil {
		glog.Error(err)
		return nil, &Error{Status: 500, Error: "Error Creating Data"}
	}

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
		glog.Error(err)
		return nil, &Error{Status: 404, Error: "This ID doesn't exist"}
	}

	sqlStatement = `
		UPDATE Shop 
		SET Name = $2, Location = $3, AdminUsername = $4, AdminPassword = $5
		WHERE id = $1;`
	_, err = db.Exec(sqlStatement, id, shop.Name, shop.Location, shop.AdminUsername, shop.AdminPassword)
	if err != nil {
		glog.Error(err)
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
		glog.Error(err)
		return nil, &Error{Status: 404, Error: "This ID doesn't exist"}
	}

	sqlStatement = `DELETE FROM Shop WHERE ID = $1;`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		glog.Error(err)
		return nil, &Error{Status: 500, Error: "Error Deleting Data"}
	}
	return &shop, nil
}

// As a Shop i can Add an item.
func ShopAddItem(id int, r *http.Request) (*Item, *Error) {
	db := ConnectToDatabase()
	defer db.Close()

	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		glog.Error(err)
		return nil, nil
	}

	item.ShopID = id
	modifiedBody, err := json.Marshal(item)
	if err != nil {
		glog.Error(err)
		return nil, nil
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(modifiedBody))
	r.ContentLength = int64(len(modifiedBody))

	addedItem, addError := CreateItem(r)
	if addError != nil {
		glog.Error(addError)
		return nil, addError
	}
	return addedItem, nil
}

// As a Shop i can Delete an item.
func ShopDeleteItem(id int, itemID int, r *http.Request) (*Item, *Error) {
	db := ConnectToDatabase()
	defer db.Close()

	item, getError := GetItem(itemID)
	if getError != nil {
		glog.Error(getError)
		return nil, getError
	}
	if item.ShopID != id {
		return nil, &Error{Status: 401, Error: "Unauthorized Access"}
	}
	deletedItem, deleteError := DeleteItem(itemID)
	if deleteError != nil {
		glog.Error(deleteError)
		return nil, deleteError
	}
	return deletedItem, nil
}

func ViewPendingOrders(id int) ([]Order, *Error) {
	var order Order
	db := ConnectToDatabase()
	defer db.Close()

	sqlStatement := `SELECT * FROM "Order" WHERE ShopID = $1 AND Delivered = false;`
	orders, err := db.Query(sqlStatement, id)
	if err != nil {
		glog.Error(err)
		return nil, nil
	}
	defer orders.Close()

	var allOrders []Order
	for orders.Next() {
		err = orders.Scan(&order.ID, &order.Delivered, &order.Price, &order.Date, &order.CustomerID, &order.ShopID)
		if err != nil {
			glog.Error(err)
			return nil, nil
		}
		allOrders = append(allOrders, order)
	}
	if allOrders == nil {
		return nil, &Error{Status: 404, Error: "No Pending Orders Exist"}
	}
	return allOrders, nil
}

func ViewDeliveredOrders(id int) ([]Order, *Error) {
	var order Order
	db := ConnectToDatabase()
	defer db.Close()

	sqlStatement := `SELECT * FROM "Order" WHERE ShopID = $1 AND Delivered = true;`
	orders, err := db.Query(sqlStatement, id)
	if err != nil {
		glog.Error(err)
		return nil, nil
	}
	defer orders.Close()

	var allOrders []Order
	for orders.Next() {
		err = orders.Scan(&order.ID, &order.Delivered, &order.Price, &order.Date, &order.CustomerID, &order.ShopID)
		if err != nil {
			glog.Error(err)
			return nil, nil
		}
		allOrders = append(allOrders, order)
	}
	if allOrders == nil {
		return nil, &Error{Status: 404, Error: "No Delivered Orders exist"}
	}
	return allOrders, nil
}

func DeliverOrder(orderID int, shopID int) *Error {
	db := ConnectToDatabase()
	defer db.Close()

	var order Order
	sqlStatement := `
		UPDATE "Order" 
		SET Delivered = true
		WHERE ID = $1 AND ShopID = $2;`
	_, err := db.Exec(sqlStatement, orderID, shopID)
	if err != nil {
		glog.Error(err)
		return &Error{Status: 400, Error: "Invalid Data"}
	}
	order.ID = orderID
	return nil
}

func ShopLogin(r *http.Request) (*Shop, *Error) {
	db := ConnectToDatabase()
	defer db.Close()

	type input struct {
		AdminUsername string `json:"adminUsername"`
		AdminPassword string `json:"adminPassword"`
	}
	var body input
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		glog.Error(err)
		return nil, nil
	}
	var shop Shop
	getShopSQL := `SELECT *
	FROM Shop
	WHERE AdminUsername = $1`
	getShoprErr := db.QueryRow(getShopSQL, body.AdminUsername).Scan(&shop.ID, &shop.Name, &shop.Location, &shop.AdminUsername, &shop.AdminPassword)
	if getShoprErr == sql.ErrNoRows {
		return nil, &Error{Status: 404, Error: "Username doesn't exist"}
	}
	if body.AdminPassword != shop.AdminPassword {
		return nil, &Error{Status: 400, Error: "Incorrect Password"}
	}
	return &shop, nil
}
