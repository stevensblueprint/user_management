package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"user_management/models"

	"gopkg.in/yaml.v2"
)

type DeleteUserRequest struct {
	Username string `json:"username"`
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request, filePath string) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userReq DeleteUserRequest
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

	// Delete the user
	delete(users.Users, userReq.Username)

	// Write the updated users.yaml file
	usersData, err = yaml.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to marshal users data", http.StatusInternalServerError)
		return
	}

	// Write the updated users.yaml file
	err = os.WriteFile(filePath, usersData, 0644)
	if err != nil {
		http.Error(w, "Failed to write to users.yaml file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User %s deleted successfully", userReq.Username)
}
