package main

import "log"

type Receive struct {
	ShopID  int    `json:"shopID"`
	OrderID int    `json:"orderID"`
	Date    string `json:"date"`
}

func CreateReceiveModel() {
	db := ConnectToDatabase()
	defer db.Close()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Receive (
		ShopID SERIAL,
		OrderID SERIAL,
		Date DATE DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT PK_Receive PRIMARY KEY (ShopID, OrderID),
		FOREIGN KEY (ShopID) REFERENCES Shop(ID) ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (OrderID) REFERENCES "Order"(ID) ON DELETE CASCADE ON UPDATE CASCADE
		);
	`)

	if err != nil {
		log.Fatal(err)
	}
}
