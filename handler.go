package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if len(port) == 0 {
		port = "8080"
	}

	r := http.NewServeMux()
	r.HandleFunc("/IncomingHTTP", Test)
	log.Println("Listening on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
	//router.HandleFunc("/create", CreateUser).Methods("POST")
	//router.HandleFunc("/read", ReadUser).Methods("GET")
	//router.HandleFunc("/update", UpdateUser).Methods("PUT")
	//router.HandleFunc("/delete", DeleteUser).Methods("DELETE")
}
