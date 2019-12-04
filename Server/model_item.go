package main

import "log"

type Item struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	ImageURL    string  `json:"imageURL"`
}

func CreateItemModel() {
	db := ConnectToDatabase()
	defer db.Close()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Item (
		ID SERIAL PRIMARY KEY,
		Name VARCHAR(255) NOT NULL,
		Type VARCHAR(255) NOT NULL,
		Price float(4)  NOT NULL,
		Description VARCHAR(255) NOT NULL,
		ImageURL VARCHAR(255) NOT NULL
		);
	`)

	if err != nil {
		log.Fatal(err)
	}
}
