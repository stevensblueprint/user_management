package handlers

import (
	"fmt"
	"net/http"
	"os"
	"user_management/models"

	"gopkg.in/yaml.v2"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request, filePath string) {
	// GET /v1/users/user?username={username}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the username from the query parameter
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Read the existing users.yaml file
	usersData, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Failed to open users.yaml file", http.StatusInternalServerError)
		return
	}

	var users models.Users
	err = yaml.Unmarshal(usersData, &users)
	if err != nil {
		http.Error(w, "Failed to parse users.yaml file", http.StatusInternalServerError)
		return
	}

	// Check if the username exists
	if _, exists := users.Users[username]; !exists {
		http.Error(w, "Username does not exist", http.StatusBadRequest)
		return
	}

	// Return the user
	user := users.Users[username]

	fmt.Fprint(w, user)
}
