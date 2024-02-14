package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"user_management/handlers"
	"user_management/middleware"

	"user_management/utils"

	"github.com/joho/godotenv"
)

var PORT = ":8080"
var BASE_URL = "/v1/users"

func main() {
	// Parse flags
	var dev bool
	var prod bool
	var PATH string

	flag.BoolVar(&dev, "dev", false, "Run the server in development mode")
	flag.BoolVar(&prod, "prod", false, "Run the server in production mode")
	flag.Parse()

	if dev {
		fmt.Println("Running in dev mode")
		PATH = "users.yaml"
		utils.ResetYAMLFile(PATH)
	}
	if prod {
		fmt.Println("Running in prod mode")
		PATH = os.Getenv("PATH")
	}
	if !dev && !prod {
		fmt.Println("Please specify a mode to run the server")
		os.Exit(1)
	}

	// Path to the users.yml file
	errEnvVariables := godotenv.Load()
	if errEnvVariables != nil {
		log.Fatal("Error: Environment variable PATH not set")
	}

	mux := http.NewServeMux()

	// Set up the routes
	mux.HandleFunc(BASE_URL+"/user", func(w http.ResponseWriter, r *http.Request) {
		// POST /v1/users/user
		if r.Method == http.MethodPost {
			handlers.AddUserHandler(w, r, PATH)
			return
		}

		// PUT /v1/users/user?username={username}
		if r.Method == http.MethodPut {
			handlers.UpdateUserHandler(w, r, PATH)
			return
		}

		// DELETE /v1/users/user?username={username}
		if r.Method == http.MethodDelete {
			handlers.DeleteUserHandler(w, r, PATH)
			return
		}

		// Return 404 for all other methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc(BASE_URL+"/all", func(w http.ResponseWriter, r *http.Request) {
		// GET /v1/users/all
		if r.Method == http.MethodGet {
			handlers.GetAllUsersHandler(w, r, PATH)
			return
		}

		// Return 404 for all other methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc(BASE_URL+"/user/enable", func(w http.ResponseWriter, r *http.Request) {
		// POST /v1/users/user/enable?username={username}
		if r.Method == http.MethodPost {
			handlers.EnableUserRequestHandler(w, r, PATH)
			return
		}

		// Return 404 for all other methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc(BASE_URL+"/user/disable", func(w http.ResponseWriter, r *http.Request) {
		// POST /v1/users/user/disable?username={username}
		if r.Method == http.MethodPost {
			handlers.DisableUserHandler(w, r, PATH)
			return
		}

		// Return 404 for all other methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Server is Healthy")
	})

	// Start the HTTP server
	wrapperMux := middleware.LoggingMiddleware(mux)

	// Ignore first char of PORT
	fmt.Printf("Server is running on port %s\n", PORT[1:])

	if err := http.ListenAndServe(PORT, wrapperMux); err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}

}
