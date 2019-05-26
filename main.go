package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	v1 "ticket/v1"

	"github.com/gorilla/mux"
)

func WrapAPIData(w http.ResponseWriter, r *http.Request, data interface{}, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	result, err := json.Marshal(map[string]interface{}{
		"Code":   code,
		"Status": message,
		"Data":   data,
	})
	if err == nil {
		log.Println(message)
		w.Write(result)
	} else {
		log.Println(fmt.Sprintf("can't wrap API data : %s", err))
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "homePage endpoint hit")
	WrapAPIData(w, r, "cobo", http.StatusOK, "successs")
	return
}

func handleRequest() {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequest()
	api := &v1.InDB{DB: conn.GetDB()}
	router := mux.NewRouter()
	router.HandleFunc("api/v1/product/book", api.BookTicket)

}
