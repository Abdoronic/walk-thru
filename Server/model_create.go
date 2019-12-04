package main

import "log"

type Create struct {
	OrderID    int `json:"orderID"`
	CustomerID int `json:"customerID"`
}

func CreateCreateModel() {
	db := ConnectToDatabase()
	defer db.Close()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS "Create" (
		OrderID SERIAL,
		CustomerID SERIAL,
		CONSTRAINT PK_Create PRIMARY KEY (OrderID, CustomerID),
		FOREIGN KEY (OrderID) REFERENCES "Order"(ID) ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (CustomerID) REFERENCES Customer(ID) ON DELETE CASCADE ON UPDATE CASCADE
		);
	`)

	if err != nil {
		log.Fatal(err)
	}
}
