package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Ticket struct {
	Name  string
	Price int
}

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

}
