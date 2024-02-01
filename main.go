package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"user_management/handlers"

	"github.com/joho/godotenv"
)

var PORT = ":3000"
var BASE_URL = "v1/users"

func main() {

	// Path to the users.yml file
	errEnvVariables := godotenv.Load()
	if errEnvVariables != nil {
		log.Fatal("Error: Environment variable PATH not set")
	}

	PATH := os.Getenv("PATH")

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
