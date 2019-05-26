package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"ticket/utils"

	"github.com/gorilla/mux"
)

type Book struct {
	Item     Product  `json:"item"`
	Quantity int      `json:"quantity"`
	Total    int      `json:"total"`
	Status   string   `json:"status"`
	Customer Customer `json:"customer"`
}

type BookRepo struct {
	Id          int    `db:"id"`
	Id_Item     int    `db:"id_item" json:"idItem"`
	Quantity    int    `db:"quantity" json:"quantity"`
	Status      string `db:"status" json:"status"`
	Total       int    `db:"total" json:"total"`
	Id_Customer int    `db:"id_customer" json:"idCustomer"`
}

func (db *InDB) TicketAvailable(id_event int, id_item int, qty int) bool {
	tx := db.DB.MustBegin()
	id := 0
	tx.Get(&id, fmt.Sprintf("SELECT id FROM products WHERE id = %d AND id_event = %d AND quantity >= %d", id_item, id_event, qty))
	if err := tx.Commit(); err != nil {
		return false
	}
	if id > 0 {
		return true
	}
	return false
}

func (db *InDB) OrderController(w http.ResponseWriter, r *http.Request) {
	id_event, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.WrapAPIError(w, r, "error converting string to integer", http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		var newBooking BookRepo
		var item Product_DB
		if err := json.NewDecoder(r.Body).Decode(&newBooking); err != nil {
			utils.WrapAPIError(w, r, "Can't decode request body", http.StatusBadRequest)
			return
		}
		if !db.TicketAvailable(id_event, newBooking.Id_Item, newBooking.Quantity) {
			utils.WrapAPIError(w, r, "Ticket is not available", http.StatusBadRequest)
			return
		}
		tx := db.DB.MustBegin()
		tx.Select(&item, fmt.Sprintf("SELECT * FROM products WHERE id = %d AND id_event = %d", newBooking.Id_Item, id_event))
		tx.MustExec("INSERT INTO bookings (id_item, id_customer, quantity, total, status) VALUES (?, ?, ?, ?, ?)", newBooking.Id_Item, newBooking.Id_Customer, newBooking.Quantity, item.Harga, "1")
		tx.MustExec("UPDATE products SET quantity = quantity - ? WHERE id = ? and quantity > 0", newBooking.Quantity, newBooking.Id_Item)
		if err := tx.Commit(); err != nil {
			utils.WrapAPIError(w, r, "error creating new booking", http.StatusInternalServerError)
			return
		}
		utils.WrapAPIData(w, r, newBooking, http.StatusOK, "success")
		return
	} else if r.Method == "GET" {
		var bookings []BookRepo
		tx := db.DB.MustBegin()
		tx.Select(&bookings, fmt.Sprintf("SELECT * FROM bookings WHERE id_event = %d", id_event))
		if err := tx.Commit(); err != nil {
			utils.WrapAPIError(w, r, "error getting ticket", http.StatusInternalServerError)
			return
		}
		response := make([]*Book, len(bookings))
		for i, item := range bookings {
			var items Product
			db.DB.Select(&item, "SELECT * FROM products WHERE id = ?", item.Id_Item)
			var customer Customer
			db.DB.Select(&customer, "SELECT * FROM customers WHERE id = ?", item.Id_Customer)
			response[i] = &Book{
				Item:     items,
				Quantity: item.Quantity,
				Total:    item.Total,
				Status:   item.Status,
				Customer: customer,
			}
		}
		utils.WrapAPIData(w, r, response, http.StatusOK, "success")
		return
	} else if r.Method == "UPDATE" {
		id_booking, err := strconv.Atoi(mux.Vars(r)["booking"])
		if err != nil {
			utils.WrapAPIError(w, r, "error converting string to integer", http.StatusInternalServerError)
			return
		}
		var booking BookRepo
		if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
			utils.WrapAPIError(w, r, "Can't decode request body", http.StatusBadRequest)
			return
		}
		tx := db.DB.MustBegin()
		tx.MustExec("UPDATE bookings SET status = ? WHERE id = ?", booking.Status, id_booking)
		if err := tx.Commit(); err != nil {
			utils.WrapAPIError(w, r, "error updating booking status", http.StatusInternalServerError)
			return
		}
		utils.WrapAPIData(w, r, booking, http.StatusOK, "success")
		return

	} else if r.Method == "DELETE" {
		id_booking, err := strconv.Atoi(mux.Vars(r)["booking"])
		if err != nil {
			utils.WrapAPIError(w, r, "error converting string to integer", http.StatusInternalServerError)
			return
		}
		var book BookRepo
		tx := db.DB.MustBegin()
		tx.Select(&book, "SELECT * FROM bookings WHERE id = ?", id_booking)
		if book.Status == "1" {
			tx.MustExec("UPDATE products SET quantity = quantity + ? WHERE id = ?", book.Quantity, book.Id_Item)
		}
		tx.MustExec("DELETE from bookings where id = ? AND id_barang = ?", id_booking, book.Id_Item)
		if err := tx.Commit(); err != nil {
			utils.WrapAPIError(w, r, "error delete booking", http.StatusInternalServerError)
			return
		}
		utils.WrapAPISuccess(w, r, "success deleting booking", http.StatusOK)
		return
	}
}
