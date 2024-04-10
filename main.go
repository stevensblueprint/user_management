package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"user_management/handlers"
	"user_management/middleware"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
)

const (
	configPath = "config.toml"
)

var PORT = ":8080"
var BASE_URL = "/api/v1/users"
var CONFIG_FILE = koanf.New(".")

func init() {
	if err := CONFIG_FILE.Load(file.Provider(configPath), toml.Parser()); err != nil {
		log.Fatalf("Error loading config file: %s", err)
	}

	if CONFIG_FILE.String("BASE_URL") == "" {
		log.Fatal("Base URL is required in config file.")
	}

	if CONFIG_FILE.String("FILE_PATH") == "" {
		log.Fatal("File path is required in config file.")
	}

	if CONFIG_FILE.String("SECRET") == "" {
		log.Fatal("Secret is required in config file.")
	}

	if CONFIG_FILE.String("redis.HOST") == "" {
		log.Fatal("Redis host is required in config file.")
	}

	if CONFIG_FILE.String("redis.PORT") == "" {
		log.Fatal("Redis port is required in config file.")
	}

	if CONFIG_FILE.String("smtp.HOST") == "" {
		log.Fatal("SMTP host is required in config file.")
	}

	if CONFIG_FILE.String("smtp.PORT") == "" {
		log.Fatal("SMTP port is required in config file.")
	}

	if CONFIG_FILE.String("smtp.USERNAME") == "" {
		log.Fatal("SMTP username is required in config file.")
	}

	if CONFIG_FILE.String("smtp.PASSWORD") == "" {
		log.Fatal("SMTP password is required in config file.")
	}
}

func main() {
	filePath := CONFIG_FILE.String("FILE_PATH")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     CONFIG_FILE.String("redis.HOST") + ":" + CONFIG_FILE.String("redis.PORT"),
		Password: "",
		DB:       0,
	})
	ctx := context.Background()

	mux := http.NewServeMux()

	// Static route
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	// Set up the routes
	mux.HandleFunc(BASE_URL+"/user", func(w http.ResponseWriter, r *http.Request) {
		// GET /v1/users/user?username={username}
		if r.Method == http.MethodGet {
			handlers.GetUserHandler(w, r, filePath)
			return
		}

		// POST /v1/users/user
		if r.Method == http.MethodPost {
			secret := CONFIG_FILE.String("SECRET")
			handlers.AddUserHandler(w, r, filePath, secret, redisClient, ctx)
			return
		}

		// PUT /v1/users/user?username={username}
		if r.Method == http.MethodPut {
			handlers.UpdateUserHandler(w, r, filePath)
			return
		}

		// DELETE /v1/users/user?username={username}
		if r.Method == http.MethodDelete {
			handlers.DeleteUserHandler(w, r, filePath)
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
			handlers.GetAllUsersHandler(w, r, filePath)
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
			handlers.EnableUserRequestHandler(w, r, filePath)
			return
		}

		// Return 405 for all other methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc(BASE_URL+"/user/disable", func(w http.ResponseWriter, r *http.Request) {
		// POST /v1/users/user/disable?username={username}
		if r.Method == http.MethodPost {
			handlers.DisableUserHandler(w, r, filePath)
			return
		}

		// Return 405 for all other methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc(BASE_URL+"/register", func(w http.ResponseWriter, r *http.Request) {
		// GET /v1/users/register?displayname={displayName}?token={token}
		if r.Method == http.MethodGet {
			handlers.RegisterPageHandler(w, r)
			return
		}

		// POST /v1/users/register
		if r.Method == http.MethodPost {
			handlers.RegisterUserHandler(w, r, CONFIG_FILE, redisClient, ctx)
			return
		}

		// Return 405 for all other methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc(BASE_URL+"/reset_password", func(w http.ResponseWriter, r *http.Request) {
		// PUT /v1/users/reset_password
		if r.Method == http.MethodPut {
			handlers.ResetPasswordHandler(w, r, filePath)
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

	corsOpt := cors.Options{
		AllowedOrigins:   []string{"admin.sitblueprint.com"},
		AllowCredentials: true,
	}

	// Start the HTTP server
	handler := cors.New(corsOpt).Handler(
		middleware.LoggingMiddleware(mux),
	)

	// Ignore first char of PORT
	fmt.Printf("Server is running on port %s\n", PORT[1:])

	if err := http.ListenAndServe(PORT, handler); err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}

}
