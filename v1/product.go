package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"ticket/utils"

	"github.com/gorilla/mux"
)

type Product struct {
	TicketName  string `json:"ticketName"`
	TicketClass string `json:"ticketClass"`
	Quantity    int    `json:"quantity"`
	Description string `json:"desc"`
	Price       int    `json:"price"`
	PictUrl     string `json:"pictUrl"`
}

type Product_DB struct {
	Id           int    `db:"id"`
	Ticket_Name  string `db:"ticket_name"`
	Ticket_Class string `db:"ticket_class"`
	Qty          int    `db:"quantity"`
	Desc         string `db:"description"`
	Harga        int    `db:"price"`
	Url_Pic      string `db:"url_pic"`
	Id_Event     int    `db:"id_event"`
}

func (db *InDB) AddProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	var id int

	if r.Method != "POST" {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.WrapAPIError(w, r, "Can't decode request body", http.StatusBadRequest)
		return
	}

	tx := db.DB.MustBegin()
	tx.Get(&id, fmt.Sprintf("SELECT id FROM EVENTS WHERE id_user = %d", id))

	tx.MustExec("INSERT INTO products (ticket_name,ticket_class,quantity,description,price,url_pic,id_event) VALUES (?, ? ,? , ?, ?, ?)", product.TicketName, product.TicketClass, product.Quantity, product.Description, product.Price, product.PictUrl, id)

	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w, r, "error adding product", http.StatusInternalServerError)
		return
	}

	utils.WrapAPISuccess(w, r, "success adding product", http.StatusCreated)
}

func (db *InDB) ListProduct(w http.ResponseWriter, r *http.Request) {

	var product []Product_DB
	var id int
	if r.Method != "GET" {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	tx := db.DB.MustBegin()
	tx.Get(&id, fmt.Sprintf("SELECT id FROM EVENTS WHERE id_user = %d", id))
	tx.Select(&product, fmt.Sprintf("SELECT * FROM PRODUCTS WHERE id_event = %d", id))

	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w, r, "error getting ticket", http.StatusInternalServerError)
		return
	}

	utils.WrapAPIData(w, r, product, http.StatusOK, "success")
}

func (db *InDB) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var id int
	var id_product int
	var err error
	if r.Method != "DELETE" {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	id_product_temp := mux.Vars(r)["id"]
	if id_product, err = strconv.Atoi(id_product_temp); err != nil {
		utils.WrapAPIError(w, r, "error converting string to integer", http.StatusInternalServerError)
		return
	}

	tx := db.DB.MustBegin()
	tx.Get(&id, fmt.Sprintf("SELECT id FROM EVENTS WHERE id_user= %d", id))
	tx.MustExec(fmt.Sprintf("DELETE from products where id_events = %d AND id = %d", id, id_product))

	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w, r, "error getting product", http.StatusInternalServerError)
		return
	}
	utils.WrapAPISuccess(w, r, "success deleting product", http.StatusOK)
}
