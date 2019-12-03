package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func GetItems() []Item {
	var item Item
	db := ConnectToDatabase()
	defer db.Close()

	sqlStatement := `SELECT * FROM Item;`
	items, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer items.Close()

	var allItems []Item
	for items.Next() {
		err = items.Scan(&item.ID, &item.Name, &item.Type, &item.Price)
		if err != nil {
			log.Fatal(err)
		}
		allItems = append(allItems, item)
	}
	return allItems
}

func GetItem(id int) (*Item, *Error) {
	var item Item
	db := ConnectToDatabase()
	defer db.Close()

	sqlStatement := `SELECT * FROM Item WHERE ID = $1;`
	err := db.QueryRow(sqlStatement, id).Scan(&item.ID, &item.Name, &item.Type, &item.Price)
	if err != nil {
		return nil, &Error{Status: 404, Error: "This ID doesn't exist"}
	}
	return &item, nil
}

func CreateItem(r *http.Request) (*Item, *Error) {
	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		return nil, &Error{Status: 400, Error: "Invalid Data"}
	}
	db := ConnectToDatabase()
	defer db.Close()

	//var id int
	sqlStatement := `
	INSERT INTO Item (Name, Type, Price)
	VALUES ($1, $2, $3) RETURNING ID, Name, Type, Price`
	err = db.QueryRow(sqlStatement, item.Name, item.Type, item.Price).Scan(&item.ID, &item.Name, &item.Type, &item.Price)
	if err != nil {
		log.Fatal(err)
		return nil, &Error{Status: 500, Error: "Error Creating Data"}
	}

	// sqlStatement = `SELECT * FROM "Item" WHERE ID = $1;`
	// _ = db.QueryRow(sqlStatement, id).Scan(&item.ID, &item.Name)

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
	err = db.QueryRow(sqlStatement, id).Scan(&temp.ID, &temp.Name, &temp.Type, &temp.Price)
	if err != nil {
		return nil, &Error{Status: 404, Error: "This ID doesn't exist"}
	}

	sqlStatement = `
		UPDATE Item 
		SET Name = $2, Type = $3, Price = $4
		WHERE id = $1;`
	_, err = db.Exec(sqlStatement, id, item.Name, item.Type, item.Price)
	if err != nil {
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
	err := db.QueryRow(sqlStatement, id).Scan(&item.ID, &item.Name, &item.Type, &item.Price)
	if err != nil {
		return nil, &Error{Status: 404, Error: "This ID doesn't exist"}
	}

	sqlStatement = `DELETE FROM Item WHERE ID = $1;`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatal(err)
		return nil, &Error{Status: 500, Error: "Error Deleting Data"}
	}
	return &item, nil
}
