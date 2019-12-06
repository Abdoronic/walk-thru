package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetShopsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetShops())
}

func GetShopHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	shop, readError := GetShop(id)
	if readError != nil {
		ErrorHandler(readError.Error, readError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(shop)
}

func CreateShopHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	shop, createError := CreateShop(r)
	if createError != nil {
		ErrorHandler(createError.Error, createError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(shop)
}

func UpdateShopHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	var updatedShop, updateError = UpdateShop(id, r)
	if updateError != nil {
		ErrorHandler(updateError.Error, updateError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(updatedShop)
}

func DeleteShopHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	var deletedShop, deleteError = DeleteShop(id)
	if deleteError != nil {
		ErrorHandler(deleteError.Error, deleteError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(deletedShop)
}

func ShopAddItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	var addedItem, addError = ShopAddItem(id, r)
	if addError != nil {
		ErrorHandler(addError.Error, addError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(addedItem)
}

func ShopDeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, errShopID := strconv.Atoi(params["id"])
	itemID, errItemID := strconv.Atoi(params["itemID"])
	if errShopID != nil || errItemID != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	var deletedItem, deleteError = ShopDeleteItem(id, itemID, r)
	if deleteError != nil {
		ErrorHandler(deleteError.Error, deleteError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(deletedItem)
}

func ViewPendingOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	json.NewEncoder(w).Encode(ViewPendingOrders(id))
}

func ViewDeliveredOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	json.NewEncoder(w).Encode(ViewDeliveredOrders(id))
}

func DeliverOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	orderID, errorderID := strconv.Atoi(params["orderID"])
	shopID, errShopID := strconv.Atoi(params["shopID"])
	if errorderID != nil || errShopID != nil {
		ErrorHandler("Invalid ID", 400, w, r)
		return
	}
	var DeliverOrderError = DeliverOrder(orderID, shopID)
	if DeliverOrderError != nil {
		ErrorHandler(DeliverOrderError.Error, DeliverOrderError.Status, w, r)
		return
	}
	json.NewEncoder(w).Encode(DeliverOrder)
}
