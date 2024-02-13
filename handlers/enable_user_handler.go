package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"user_management/models"

	"gopkg.in/yaml.v2"
)

type EnableUserRequest struct {
	Username string `json:"username"`
}

func EnableUserRequestHandler(w http.ResponseWriter, r *http.Request, filePath string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userReq EnableUserRequest
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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
	if _, exists := users.Users[userReq.Username]; !exists {
		http.Error(w, "Username does not exist", http.StatusBadRequest)
		return
	}

	// Check if user is already enabled
	if !users.Users[userReq.Username].Disabled {
		http.Error(w, "User is already enabled", http.StatusBadRequest)
		return
	}

	// Enable the user
	user := users.Users[userReq.Username]
	user.Disabled = false
	users.Users[userReq.Username] = user

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

	fmt.Fprint(w, "User enabled successfully")
}
