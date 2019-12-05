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
	DBHost     string `json:"dbHost"`
	DBPort     string `json:"dbPort"`
	DBUser     string `json:"dbUser"`
	DBPassword string `json:"dbPassword"`
	DBName     string `json:"dbName"`
	WebHost    string `json:"webHost"`
	WebPort    string `json:"webPort"`
}

func main() {

	cfg := GetConfig()

	CreateModels()

	router := CreateRouter()

	fmt.Println("Listening on " + cfg.WebHost + "Port: " + cfg.WebPort)
	log.Fatal(http.ListenAndServe(cfg.WebHost+":"+cfg.WebPort, router))
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
		cfg.DBHost = os.Getenv("DATABASE_HOST")
		cfg.DBName = os.Getenv("DATABASE_DBNAME")
		cfg.DBPort = os.Getenv("DATABASE_PORT")
		cfg.DBUser = os.Getenv("DATABASE_USER")
		cfg.DBPassword = os.Getenv("DATABASE_PASSWORD")
		cfg.WebHost = os.Getenv("WEB_HOST")
		cfg.WebPort = os.Getenv("WEB_PORT")
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
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

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
