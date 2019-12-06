package main

import "log"

type Contain struct {
	OrderID  int `json:"orderID"`
	ItemID   int `json:"itemID"`
	Quantity int `json:"quantity"`
}

func CreateContainModel() {
	db := ConnectToDatabase()
	defer db.Close()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Contain (
		OrderID SERIAL,
		ItemID SERIAL,
		Quantity INT,
		CONSTRAINT PK_Contain PRIMARY KEY (OrderID, ItemID),
		FOREIGN KEY (OrderID) REFERENCES "Order"(ID) ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (ItemID) REFERENCES Item(ID) ON DELETE CASCADE ON UPDATE CASCADE
		);
	`)

	if err != nil {
		log.Fatal(err)
	}
}
