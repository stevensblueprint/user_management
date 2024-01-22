package main

import (
	"fmt"
	"log"
	"net/http"

	"user_management/handlers"
)

var PORT = ":3000"
var PATH = "etc/users.yaml"

func main() {
	// Set up the routes
	http.HandleFunc("/add-user", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddUserHandler(w, r, PATH)
	})
	http.Handle("/", http.FileServer(http.Dir("static")))

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Server is Healthy")
	})

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(PORT, nil))
}
