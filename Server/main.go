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
	Port     string `json:"port"`
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
	CreateOfferModel()
	CreateContainModel()
	CreateReceiveModel()
	CreateCreateModel()
}

func GetConfig() *Config {
	var cfg Config
	if value := os.Getenv("DATABASE_HOST"); value != "" {
		cfg.Host = os.Getenv("DATABASE_HOST")
		cfg.DbName = os.Getenv("DATABASE_DBNAME")
		cfg.Port = os.Getenv("DATABASE_PORT")
		cfg.User = os.Getenv("DATABASE_USER")
		cfg.Password = os.Getenv("DATABASE_PASSWORD")
	} else {
		jsonFile, err := os.Open("config.json")
		if err != nil {
			log.Fatal(err)
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &cfg)
	}
	return &cfg
}

func ConnectToDatabase() *sql.DB {
	cfg := GetConfig()
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
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
