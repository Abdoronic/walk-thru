package main

import (
	"encoding/json"
	"net/http"

	"github.com/golang/glog"
)

func GetItems() ([]Item, *Error) {
	var item Item
	db := ConnectToDatabase()
	defer db.Close()

	sqlStatement := `SELECT * FROM Item;`
	items, err := db.Query(sqlStatement)
	if err != nil {
		glog.Error(err)
		return nil, nil
	}
	defer items.Close()

	var allItems []Item
	for items.Next() {
		err = items.Scan(&item.ID, &item.Name, &item.Type, &item.Price, &item.Description, &item.ImageURL, &item.ShopID)
		if err != nil {
			glog.Error(err)
			return nil, nil
		}
		allItems = append(allItems, item)
	}
	if allItems == nil {
		return nil, &Error{Status: 404, Error: "No Items exists"}
	}
	return allItems, nil
}

func GetItem(id int) (*Item, *Error) {
	var item Item
	db := ConnectToDatabase()
	defer db.Close()

	sqlStatement := `SELECT * FROM Item WHERE ID = $1;`
	err := db.QueryRow(sqlStatement, id).Scan(&item.ID, &item.Name, &item.Type, &item.Price, &item.Description, &item.ImageURL, &item.ShopID)
	if err != nil {
		glog.Error(err)
		return nil, &Error{Status: 404, Error: "This ID doesn't exist"}
	}
	return &item, nil
}

func CreateItem(r *http.Request) (*Item, *Error) {
	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		glog.Error(err)
		return nil, &Error{Status: 400, Error: "Invalid Data"}
	}
	db := ConnectToDatabase()
	defer db.Close()

	sqlStatement := `
	INSERT INTO Item (Name, Type, Price, Description, ImageURL, ShopID)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING ID, Name, Type, Price, Description, ImageURL, ShopID`
	err = db.QueryRow(sqlStatement, item.Name, item.Type, item.Price, item.Description, item.ImageURL, item.ShopID).Scan(&item.ID, &item.Name, &item.Type, &item.Price, &item.Description, &item.ImageURL, &item.ShopID)
	if err != nil {
		glog.Error(err)
		return nil, &Error{Status: 500, Error: "Error Creating Data"}
	}

	return &item, nil
}

func UpdateItem(id int, r *http.Request) (*Item, *Error) {
	db := ConnectToDatabase()
	defer db.Close()

	var (
		item Item
		temp Item
	)
	err := json.NewDecoder(r.Body).Decode(&item)

	sqlStatement := `SELECT * FROM Item WHERE ID = $1;`
	err = db.QueryRow(sqlStatement, id).Scan(&temp.ID, &temp.Name, &temp.Type, &temp.Price, &temp.Description, &temp.ImageURL, &temp.ShopID)
	if err != nil {
		glog.Error(err)
		return nil, &Error{Status: 404, Error: "This ID doesn't exist"}
	}

	sqlStatement = `
		UPDATE Item 
		SET Name = $2, Type = $3, Price = $4, Description = $5, ImageURL = $6, ShopID = $7
		WHERE id = $1;`
	_, err = db.Exec(sqlStatement, id, item.Name, item.Type, item.Price, item.Description, item.ImageURL, item.ShopID)
	if err != nil {
		glog.Error(err)
		return nil, &Error{Status: 400, Error: "Invalid Data"}
	}
	item.ID = id
	return &item, nil
}

func DeleteItem(id int) (*Item, *Error) {
	db := ConnectToDatabase()
	defer db.Close()

	var item Item
	sqlStatement := `SELECT * FROM Item WHERE ID = $1;`
	err := db.QueryRow(sqlStatement, id).Scan(&item.ID, &item.Name, &item.Type, &item.Price, &item.Description, &item.ImageURL, &item.ShopID)
	if err != nil {
		glog.Error(err)
		return nil, &Error{Status: 404, Error: "This ID doesn't exist"}
	}

	sqlStatement = `DELETE FROM Item WHERE ID = $1;`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		glog.Error(err)
		return nil, &Error{Status: 500, Error: "Error Deleting Data"}
	}
	return &item, nil
}
