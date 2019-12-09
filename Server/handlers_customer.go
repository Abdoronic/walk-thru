package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetCustomersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetCustomers())
}

func GetCustomerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	customer, readError := GetCustomer(id)
	if readError != nil {
		ErrorHandler(readError.Error, readError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(customer)
}

func CreateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	customer, createError := CreateCustomer(r)
	if createError != nil {
		ErrorHandler(createError.Error, createError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(customer)
}

func UpdateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	var updatedCustomer, updateError = UpdateCustomer(id, r)
	if updateError != nil {
		ErrorHandler(updateError.Error, updateError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(updatedCustomer)
}

func DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	var deletedCustomer, deleteError = DeleteCustomer(id)
	if deleteError != nil {
		ErrorHandler(deleteError.Error, deleteError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(deletedCustomer)
}

func ViewOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	json.NewEncoder(w).Encode(ViewCustomerOrders(id))
}

func ViewItemsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	json.NewEncoder(w).Encode(ViewItems(id))
}

func CustomerCreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	var createdOrder, createError = CustomerCreateOrder(id, r)
	if createError != nil {
		ErrorHandler(createError.Error, createError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(createdOrder)
}

func ViewOrderItemsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	orderID, err := strconv.Atoi(params["orderID"])
	if err != nil {
		ErrorHandler("Invalid Order ID", 400, w, r)
		return
	}
	items, viewError := ViewOrderItems(orderID)
	if viewError != nil {
		ErrorHandler(viewError.Error, viewError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(items)
}

func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var (
		customerID int
		orderID    int
		shopID     int
		err        error
	)
	customerID, err = strconv.Atoi(params["customerID"])
	if err != nil {
		ErrorHandler("Invalid Customer ID", 400, w, r)
		return
	}
	orderID, err = strconv.Atoi(params["orderID"])
	if err != nil {
		ErrorHandler("Invalid Order ID", 400, w, r)
		return
	}
	shopID, err = strconv.Atoi(params["shopID"])
	if err != nil {
		ErrorHandler("Invalid Shop ID", 400, w, r)
		return
	}
	checkedOrder, checkoutError := Checkout(customerID, orderID, shopID, r)
	if err != nil {
		ErrorHandler(checkoutError.Error, checkoutError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(checkedOrder)
}

func CustomerAddItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	orderID, orderErr := strconv.Atoi(params["orderID"])
	itemID, itemErr := strconv.Atoi(params["itemID"])
	if orderErr != nil || itemErr != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	var updatedOrder, updateError = CustomerAddItem(orderID, itemID, r)
	if updateError != nil {
		ErrorHandler(updateError.Error, updateError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(updatedOrder)
}

func CustomerRemoveItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	orderID, orderErr := strconv.Atoi(params["orderID"])
	itemID, itemErr := strconv.Atoi(params["itemID"])
	if orderErr != nil || itemErr != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	var updatedOrder, updateError = CustomerRemoveItem(orderID, itemID, r)
	if updateError != nil {
		ErrorHandler(updateError.Error, updateError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(updatedOrder)
}

func CustomerViewOrderItemsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	customerID, customerErr := strconv.Atoi(params["id"])
	orderID, orderErr := strconv.Atoi(params["orderID"])
	if orderErr != nil || customerErr != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	var orderItems, getError = CustomerViewOrderItems(customerID, orderID, r)
	if getError != nil {
		ErrorHandler(getError.Error, getError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(orderItems)
}
func CustomerLoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var customer, getError = CustomerLogin(r)
	if getError != nil {
		ErrorHandler(getError.Error, getError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(customer)
}
