package main

import (
	"log"
	"net/http"

	"user_management/handlers"
)

func main() {
	// Add the user.yaml file path
	path := "etc/users.yaml"

	// Set up the routes
	http.HandleFunc("/add-user", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddUserHandler(w, r, path)
	})
	http.Handle("/", http.FileServer(http.Dir("static")))

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":3000", nil))
}
