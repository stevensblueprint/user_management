package handlers

import (
	"fmt"
	"net/http"
	"os"
	"user_management/models"

	"gopkg.in/yaml.v2"
)

func DisableUserHandler(w http.ResponseWriter, r *http.Request, filePath string) {
	// POST /v1/users/disable?username={username}
	if r.Method != http.MethodPost {
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

	// Check if user is already disabled
	if users.Users[username].Disabled {
		http.Error(w, "User is already disabled", http.StatusBadRequest)
		return
	}

	// Disable the user
	user := users.Users[username]
	user.Disabled = true
	users.Users[username] = user

	// Write the updated users.yaml file
	usersData, err = yaml.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to convert users to YAML", http.StatusInternalServerError)
		return
	}

	err = os.WriteFile(filePath, usersData, 0644)
	if err != nil {
		http.Error(w, "Failed to write users.yaml file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User disabled successfully")
}
