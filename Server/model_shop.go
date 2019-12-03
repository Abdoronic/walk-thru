package main

import "log"

type Shop struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Location      string `json:"location"`
	AdminUsername string `json:"adminUsername"`
	AdminPassword string `json:"adminPassword"`
}

func CreateShopModel() {
	db := ConnectToDatabase()
	defer db.Close()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Shop (
		ID SERIAL PRIMARY KEY,
		Name VARCHAR(255) NOT NULL,
		Location VARCHAR(255) NOT NULL,
		AdminUsername VARCHAR(255) UNIQUE NOT NULL,
		AdminPassword VARCHAR(255) NOT NULL
		);
	`)

	if err != nil {
		log.Fatal(err)
	}
}
