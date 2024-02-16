package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"user_management/models"

	"gopkg.in/yaml.v2"
)

type ResetPasswordRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request, filePath string) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var resetReq ResetPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&resetReq)
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

	// Check if username is empty
	if resetReq.Username == "" {
		http.Error(w, "Username cannot be empty", http.StatusBadRequest)
		return
	}

	// Check if username exists
	if _, exists := users.Users[resetReq.Username]; !exists {
		http.Error(w, "Username does not exist", http.StatusBadRequest)
		return
	}

	// Check if password is empty
	if resetReq.Password == "" {
		http.Error(w, "Password cannot be empty", http.StatusBadRequest)
		return
	}

	// Check if password is the same as the old password
	if users.Users[resetReq.Username].Password == resetReq.Password {
		http.Error(w, "Password is the same as old password", http.StatusBadRequest)
		return
	}

	// Update the password
	user := users.Users[resetReq.Username]
	user.Password = resetReq.Password

	// Convert the users object to YAML
	usersData, err = yaml.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to convert users object to YAML", http.StatusInternalServerError)
		return
	}

	// Write the updated users.yaml file
	err = os.WriteFile(filePath, usersData, 0644)
	if err != nil {
		http.Error(w, "Failed to write to users.yaml file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password updated successfully"))
}
