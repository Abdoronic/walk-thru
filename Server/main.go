package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"dbName"`
}

func main() {

	CreateModels()

	router := CreateRouter()

	fmt.Println("Listening on Port: 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func CreateModels() {
	CreateCustomerModel()
	CreateItemModel()
  CreateShopModel()
	CreateOrderModel()
}

func ConnectToDatabase() *sql.DB {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var cfg Config
	json.Unmarshal(byteValue, &cfg)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
