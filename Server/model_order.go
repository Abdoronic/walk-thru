package main

import (
	"database/sql"
	"log"
)

type Order struct {
	ID         int           `json:"id"`
	Delivered  bool          `json:"delivered"`
	Price      float64       `json:"price"`
	Date       string        `json:"date"`
	CustomerID int           `json:"customerID"`
	ShopID     sql.NullInt64 `json:"shopID"`
}

func CreateOrderModel() {
	db := ConnectToDatabase()
	defer db.Close()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS "Order" (
		ID SERIAL PRIMARY KEY,
		Delivered BOOLEAN DEFAULT FALSE,
		Price float(4)  NOT NULL,
		Date DATE DEFAULT CURRENT_TIMESTAMP,
		CustomerID SERIAL,
		ShopID SERIAL,
		FOREIGN KEY (CustomerID) REFERENCES Customer(ID) ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (ShopID) REFERENCES Shop(ID) ON DELETE CASCADE ON UPDATE CASCADE
		);
	`)

	if err != nil {
		log.Fatal(err)
	}
}
