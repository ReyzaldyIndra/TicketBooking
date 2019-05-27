package main

import (
	"fmt"
	"log"
	"net/http"
	"ticket/database"
	v1 "ticket/v1"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// func WrapAPIData(w http.ResponseWriter, r *http.Request, data interface{}, code int, message string) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	result, err := json.Marshal(map[string]interface{}{
// 		"Code":   code,
// 		"Status": message,
// 		"Data":   data,
// 	})
// 	if err == nil {
// 		log.Println(message)
// 		w.Write(result)
// 	} else {
// 		log.Println(fmt.Sprintf("can't wrap API data : %s", err))
// 	}
// }

// func homePage(w http.ResponseWriter, r *http.Request) {
// 	// fmt.Fprintf(w, "homePage endpoint hit")
// 	WrapAPIData(w, r, "cobo", http.StatusOK, "successs")
// 	return
// }

// func handleRequest() {
// 	http.HandleFunc("/", homePage)
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

func main() {
	// handleRequest()
	viper.SetConfigFile("./config/dev.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	conn, err := database.InitDb(fmt.Sprintf("%s:%s@tcp(%s:%s)/commerce", viper.Get("db.username"), viper.Get("db.password"), viper.Get("db.host"), viper.Get("db.port")))
	if err != nil {
		fmt.Errorf("failed to open database: %v", err)
		return
	}
	defer conn.DB.Close()

	api := &v1.InDB{DB: conn.GetDB()}
	router := mux.NewRouter()
	//API Mapping
	router.HandleFunc("/api/v1/event/{id}/book/{book}", http.HandlerFunc(api.OrderController))
	router.HandleFunc("/api/v1/event/{id}/user/{user}", http.HandlerFunc(api.CustomerController))
	router.HandleFunc("api/v1/event/create", http.HandlerFunc(api.CreateEvent))
	router.HandleFunc("/api/v1/product/add", http.HandlerFunc(api.AddProduct))
	router.HandleFunc("/api/v1/product/list", http.HandlerFunc(api.ListProduct))
	router.HandleFunc("/api/v1/product/delete/{id}", http.HandlerFunc(api.DeleteProduct))
	router.HandleFunc("/api/v1/user/create", api.CreateUser)

	http.Handle("/", router)
	port := fmt.Sprintf(":%s", viper.Get("host.port"))
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
