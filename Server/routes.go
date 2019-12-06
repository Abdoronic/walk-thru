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

	router.HandleFunc("/customers/view/shops", GetShopsHandler).Methods("GET")
	router.HandleFunc("/customers/viewItems/{id}", ViewItemsHandler).Methods("GET")
	router.HandleFunc("/customers/viewOrderItems/{orderID}", ViewOrderItemsHandler).Methods("GET")
	router.HandleFunc("/customers/{id}/createOrder", CustomerCreateOrderHandler).Methods("POST")
	router.HandleFunc("/customers/{id}/addItem/{orderID}/{itemID}", CustomerAddItemHandler).Methods("PUT")
	router.HandleFunc("/customers/{id}/removeItem/{orderID}/{itemID}", CustomerRemoveItemHandler).Methods("PUT")
	router.HandleFunc("/customers/{customerID}/checkout/{orderID}/{shopID}", CheckoutHandler).Methods("POST")
	router.HandleFunc("/customers/{id}/viewOrderItems/{orderID}", CustomerViewOrderItemsHandler).Methods("GET")
	router.HandleFunc("/customers/login", CustomerLoginHandler).Methods("POST")

	router.HandleFunc("/items", GetItemsHandler).Methods("GET")
	router.HandleFunc("/items/{id}", GetItemHandler).Methods("GET")
	router.HandleFunc("/items", CreateItemHandler).Methods("POST")
	router.HandleFunc("/items/{id}", UpdateItemHandler).Methods("PUT")
	router.HandleFunc("/items/{id}", DeleteItemHandler).Methods("DELETE")

	router.HandleFunc("/shops/{id}/viewPendingOrders", ViewPendingOrdersHandler).Methods("GET")
	router.HandleFunc("/shops/{id}/viewDeliveredOrders", ViewDeliveredOrdersHandler).Methods("GET")
	router.HandleFunc("/shops/{id}/viewOfferedItems", ViewItemsHandler).Methods("GET")
	router.HandleFunc("/shops/{shopID}/deliverOrder/{orderID}", DeliverOrderHandler).Methods("PUT")

	router.HandleFunc("/shops", GetShopsHandler).Methods("GET")
	router.HandleFunc("/shops/{id}", GetShopHandler).Methods("GET")
	router.HandleFunc("/shops", CreateShopHandler).Methods("POST")
	router.HandleFunc("/shops/{id}", UpdateShopHandler).Methods("PUT")
	router.HandleFunc("/shops/{id}", DeleteShopHandler).Methods("DELETE")

	router.HandleFunc("/shops/{id}/shopAddItem", ShopAddItemHandler).Methods("POST")
	router.HandleFunc("/shops/{id}/shopDeleteItem/{itemID}", ShopDeleteItemHandler).Methods("DELETE")
	router.HandleFunc("/shops/login", ShopLoginHandler).Methods("POST")

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
