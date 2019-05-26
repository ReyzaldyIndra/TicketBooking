package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"ticket/utils"

	"github.com/gorilla/mux"
)

type Customer struct {
	Email     string `json:"email"`
	Nama      string `json:"nama"`
	Destinasi string `json:"destination"`
}

func (db *InDB) isMyEvent(id_user int, id_store int) bool {
	id := 0
	tx := db.DB.MustBegin()
	tx.Get(&id, fmt.Sprintf("SELECT id FROM events WHERE id = %d AND id_user = %d", id_store, id_user))
	if err := tx.Commit(); err != nil {
		return false
	}
	if id > 0 {
		return true
	}
	return false
}

func (db *InDB) CustomerController(w http.ResponseWriter, r *http.Request) {
	id_event, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.WrapAPIError(w, r, "error converting string to integer", http.StatusInternalServerError)
		return
	}
	var id_user int
	if !db.isMyEvent(id_user, id_event) {
		utils.WrapAPIError(w, r, "invalid token", http.StatusBadRequest)
		return
	}

	if r.Method == "POST" {
		var customer Customer
		if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
			utils.WrapAPIError(w, r, "Can't decode request body", http.StatusBadRequest)
			return
		}
		id := 0
		tx := db.DB.MustBegin()
		tx.Get(&id, fmt.Sprintf("SELECT id FROM customers WHERE cust_email = %s AND id_event = %d", customer.Email, id_event))
		if id > 0 {
			tx.Select(&customer, fmt.Sprintf("SELECT * FROM customers WHERE id =  %d", id))
			utils.WrapAPIData(w, r, customer, http.StatusOK, "success")
			return
		}
		tx.MustExec("INSERT INTO customers (cust_name, cust_address, cust_email, id_store) VALUES (?, ?, ?, ?", customer.Nama, customer.Destinasi, customer.Email, id_event)

		if err := tx.Commit(); err != nil {
			utils.WrapAPIError(w, r, "error inserting new customer", http.StatusInternalServerError)
			return
		}
		utils.WrapAPIData(w, r, customer, http.StatusOK, "success")
		return
	} else if r.Method == "GET" {
		var customer Customer
		id_cust, err := strconv.Atoi(mux.Vars(r)["user"])
		if err != nil {
			utils.WrapAPIError(w, r, "error converting string to integer", http.StatusInternalServerError)
			return
		}
		tx := db.DB.MustBegin()
		tx.Select(&customer, fmt.Sprintf("SELECT * FROM customers WHERE id =  %d", id_cust))
		if err := tx.Commit(); err != nil {
			utils.WrapAPIError(w, r, "error get user", http.StatusInternalServerError)
			return
		}
		utils.WrapAPIData(w, r, customer, http.StatusOK, "success")
		return
	} else if r.Method == "PUT" {
		id_cust, err := strconv.Atoi(mux.Vars(r)["user"])
		if err != nil {
			utils.WrapAPIError(w, r, "error converting string to integer", http.StatusInternalServerError)
			return
		}
		var customer Customer
		if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
			utils.WrapAPIError(w, r, "Can't decode request body", http.StatusBadRequest)
			return
		}
		tx := db.DB.MustBegin()
		tx.MustExec("UPDATE customers SET cust_name = ?, cust_address = ?, cust_email = ? WHERE id = ?", customer.Nama, customer.Destinasi, customer.Email, id_cust)
		if err := tx.Commit(); err != nil {
			utils.WrapAPIError(w, r, "error updating user", http.StatusInternalServerError)
			return
		}
		utils.WrapAPIData(w, r, customer, http.StatusOK, "success")
		return
	} else if r.Method == "DELETE" {
		id_cust, err := strconv.Atoi(mux.Vars(r)["user"])
		if err != nil {
			utils.WrapAPIError(w, r, "error converting string to integer", http.StatusInternalServerError)
			return
		}
		tx := db.DB.MustBegin()
		tx.MustExec("DELETE from customers where id = ? AND id_store = %d", id_cust, id_event)
		if err := tx.Commit(); err != nil {
			utils.WrapAPIError(w, r, "error delete product", http.StatusInternalServerError)
			return
		}
		utils.WrapAPISuccess(w, r, "success deleting user", http.StatusOK)
		return
	}
	utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	return
}
