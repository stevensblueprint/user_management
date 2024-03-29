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

	"github.com/rs/cors"
)

var PORT = ":8080"
var BASE_URL = "/api/v1/users"

func main() {
	// Parse flags
	var dev bool
	var prod bool
	var help bool
	var PATH string

	flag.BoolVar(&dev, "dev", false, "Run the server in development mode")
	flag.BoolVar(&prod, "prod", false, "Run the server in production mode")
	flag.BoolVar(&help, "h", false, "Show help")
	flag.Parse()

	// Check for flags
	switch {
	case dev && prod:
		fmt.Println("Usage: Cannot run in both dev and prod mode")
		os.Exit(1)
	case !dev && !prod && !help:
		fmt.Println("Usage: main.go [-dev] [-prod] [-path PATH]")
		os.Exit(1)
	case dev:
		fmt.Println("Running in dev mode")
		PATH = "users.yaml"
		utils.ResetYAMLFile(PATH)
	case prod:
		fmt.Println("Running in prod mode")
		PATH = os.Getenv("PATH")
		if PATH == "" {
			log.Fatal("Error: Environment variable PATH not set")
		}
	case help:
		fmt.Println("Usage: main.go [-dev] [-prod] [-path PATH]")
		os.Exit(0)
	}

	mux := http.NewServeMux()

	// Set up the routes
	mux.HandleFunc(BASE_URL+"/user", func(w http.ResponseWriter, r *http.Request) {
		// GET /v1/users/user?username={username}
		if r.Method == http.MethodGet {
			handlers.GetUserHandler(w, r, PATH)
			return
		}

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

		if r.Method == http.MethodHead {
			return
		}

		// Return 405 for all other methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc(BASE_URL+"/all", func(w http.ResponseWriter, r *http.Request) {
		// GET /v1/users/all
		if r.Method == http.MethodGet {
			handlers.GetAllUsersHandler(w, r, PATH)
			return
		}

		if r.Method == http.MethodHead {
			return
		}

		// Return 405 for all other methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc(BASE_URL+"/user/enable", func(w http.ResponseWriter, r *http.Request) {
		// POST /v1/users/user/enable?username={username}
		if r.Method == http.MethodPost {
			handlers.EnableUserRequestHandler(w, r, PATH)
			return
		}

		// Return 405 for all other methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc(BASE_URL+"/user/disable", func(w http.ResponseWriter, r *http.Request) {
		// POST /v1/users/user/disable?username={username}
		if r.Method == http.MethodPost {
			handlers.DisableUserHandler(w, r, PATH)
			return
		}

		// Return 405 for all other methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc(BASE_URL+"/register", func(w http.ResponseWriter, r *http.Request) {
		// POST /v1/users/register
		if r.Method == http.MethodPost {
			handlers.RegisterUserHandler(w, r)
			return
		}

		// Return 405 for all other methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc(BASE_URL+"/reset_password", func(w http.ResponseWriter, r *http.Request) {
		// PUT /v1/users/reset_password
		if r.Method == http.MethodPut {
			handlers.ResetPasswordHandler(w, r, PATH)
			return
		}

		// Return 405 for all other methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Server is Healthy")
	})

	// Start the HTTP server
	handler := cors.Default().Handler(
		middleware.LoggingMiddleware(mux),
	)

	// Ignore first char of PORT
	fmt.Printf("Server is running on port %s\n", PORT[1:])

	if err := http.ListenAndServe(PORT, handler); err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}

}
