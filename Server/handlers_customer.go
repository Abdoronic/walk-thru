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
