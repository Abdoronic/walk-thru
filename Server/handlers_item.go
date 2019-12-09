package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetItemsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	items, readError := GetItems()
	if readError != nil {
		ErrorHandler(readError.Error, readError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(items)
}

func GetItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	item, readError := GetItem(id)
	if readError != nil {
		ErrorHandler(readError.Error, readError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(item)
}

func CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	item, createError := CreateItem(r)
	if createError != nil {
		ErrorHandler(createError.Error, createError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(item)
}

func UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	var updatedItem, updateError = UpdateItem(id, r)
	if updateError != nil {
		ErrorHandler(updateError.Error, updateError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(updatedItem)
}

func DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	var deletedItem, deleteError = DeleteItem(id)
	if deleteError != nil {
		ErrorHandler(deleteError.Error, deleteError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(deletedItem)
}
