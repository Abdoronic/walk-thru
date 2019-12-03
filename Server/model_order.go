package main

import "log"

type Order struct {
	ID        int     `json:"id"`
	Delivered bool    `json:"Delivered"`
	Price     float64 `json:"price"`
}

func CreateOrderModel() {
	db := ConnectToDatabase()
	defer db.Close()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS "Order" (
		ID SERIAL PRIMARY KEY,
		Delivered BOOLEAN DEFAULT FALSE,
		Price float(4)  NOT NULL
		);
	`)

	if err != nil {
		log.Fatal(err)
	}
}
