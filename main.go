package main

import (
	"log"
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()
	router := route.PathPrefix("/api").Subrouter()
	//Routes
	// for sequential emulation
	runtime.GOMAXPROCS(2)
	router.HandleFunc("/test", Test).Methods("GET")
	router.HandleFunc("/create", CreateUser).Methods("GET")
	router.HandleFunc("/read", ReadUser).Methods("GET")
	router.HandleFunc("/update", UpdateUser).Methods("GET")
	router.HandleFunc("/delete", DeleteUser).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router)) // Run Server
}
