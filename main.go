package main

import (
	"fmt"
	"log"
	"net/http"

	"user_management/handlers"
)

var PORT = ":3000"
var PATH = "etc/users.yaml"
var BASE_URL = "v1/users"

func main() {
	// Set up the routes
	http.HandleFunc(BASE_URL+"/user", func(w http.ResponseWriter, r *http.Request) {
		// POST /v1/users/user
		if r.Method == http.MethodPost {
			handlers.AddUserHandler(w, r, PATH)
		}

		// Return 404 for all other methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc(BASE_URL+"/user/disable", func(w http.ResponseWriter, r *http.Request) {
		// POST /v1/users/user/disable
		if r.Method == http.MethodPost {
			handlers.DisableUserHandler(w, r, PATH)
		}

		// Return 404 for all other methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc(BASE_URL+"/all", func(w http.ResponseWriter, r *http.Request) {
		// GET /v1/users/all
		if r.Method == http.MethodGet {
			handlers.GetAllUsersHandler(w, r, PATH)
		}

		// Return 404 for all other methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Server is Healthy")
	})

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(PORT, nil))
}
