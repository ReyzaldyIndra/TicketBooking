package v1

import (
	"encoding/json"
	"net/http"
	"ticket/utils"
)

type Event struct {
	EventName   string `json:"eventName"`
	Venue       string `json:"venue"`
	PhoneNumber string `json:"phone"`
	BankNumber  string `json:"bankNumber"`
}

func (db *InDB) CreateEvent(w http.ResponseWriter, r *http.Request) {

	var event Event
	var id int
	var err error
	if r.Method != "POST" {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		utils.WrapAPIError(w, r, "Can't decode request body", http.StatusBadRequest)
		return
	}

	tx := db.DB.MustBegin()
	tx.MustExec("INSERT INTO events (event_name,venue,handphone,bank_number,id_user) VALUES (?, ? ,? , ?, ?)", event.EventName, event.Venue, event.PhoneNumber, event.BankNumber, id)

	if err = tx.Commit(); err != nil {
		utils.WrapAPIError(w, r, "error creating new event", http.StatusInternalServerError)
		return
	}

	utils.WrapAPISuccess(w, r, "success creating new event", http.StatusCreated)
}
