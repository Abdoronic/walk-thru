package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/customers", GetCustomersHandler).Methods("GET")
	router.HandleFunc("/customers/{id}", GetCustomerHandler).Methods("GET")
	router.HandleFunc("/customers", CreateCustomerHandler).Methods("POST")
	router.HandleFunc("/customers/{id}", UpdateCustomerHandler).Methods("PUT")
	router.HandleFunc("/customers/{id}", DeleteCustomerHandler).Methods("DELETE")

	router.HandleFunc("/customers/view/shops", ViewShopsHandler).Methods("GET")

	router.HandleFunc("/items", GetItemsHandler).Methods("GET")
	router.HandleFunc("/items/{id}", GetItemHandler).Methods("GET")
	router.HandleFunc("/items", CreateItemHandler).Methods("POST")
	router.HandleFunc("/items/{id}", UpdateItemHandler).Methods("PUT")
	router.HandleFunc("/items/{id}", DeleteItemHandler).Methods("DELETE")

	router.HandleFunc("/shops", GetShopsHandler).Methods("GET")
	router.HandleFunc("/shops/{id}", GetShopHandler).Methods("GET")
	router.HandleFunc("/shops", CreateShopHandler).Methods("POST")
	router.HandleFunc("/shops/{id}", UpdateShopHandler).Methods("PUT")
	router.HandleFunc("/shops/{id}", DeleteShopHandler).Methods("DELETE")

	router.HandleFunc("/shops/{id}/shopAddItem", ShopAddItemHandler).Methods("POST")
	router.HandleFunc("/shops/{id}/shopDeleteItem/{itemID}", ShopDeleteItemHandler).Methods("DELETE")

	router.HandleFunc("/orders", GetOrdersHandler).Methods("GET")
	router.HandleFunc("/orders/{id}", GetOrderHandler).Methods("GET")
	router.HandleFunc("/orders", CreateOrderHandler).Methods("POST")
	router.HandleFunc("/orders/{id}", UpdateOrderHandler).Methods("PUT")
	router.HandleFunc("/orders/{id}", DeleteOrderHandler).Methods("DELETE")

	return router
}

func ErrorHandler(msg string, status int, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&Error{Status: status, Error: msg})
}

type Error struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}
