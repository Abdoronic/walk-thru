package main

import "log"

type Offer struct {
	ShopID int `json:"shopID"`
	ItemID int `json:"itemID"`
}

func CreateOfferModel() {
	db := ConnectToDatabase()
	defer db.Close()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Offer (
		ShopID SERIAL,
		ItemID SERIAL,
		CONSTRAINT PK_Offer PRIMARY KEY (ShopID, ItemID),
		FOREIGN KEY (ShopID) REFERENCES Shop(ID) ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (ItemID) REFERENCES Item(ID) ON DELETE CASCADE ON UPDATE CASCADE
		);
	`)

	if err != nil {
		log.Fatal(err)
	}
}
