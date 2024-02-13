package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"user_management/models"

	"gopkg.in/yaml.v2"
)

type UpdateUserRequest struct {
	Username    string   `json:"username"`
	Displayname string   `json:"displayname,omitempty"`
	Email       string   `json:"email,omitempty"`
	Groups      []string `json:"groups,omitempty"`
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request, filePath string) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userReq UpdateUserRequest
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

	// If email is updated, check if it ends with @sitblueprint.com or @stevens.edu
	if userReq.Email != "" {
		if !strings.HasSuffix(userReq.Email, "@sitblueprint.com") && !strings.HasSuffix(userReq.Email, "@stevens.edu") {
			http.Error(w, "Invalid email domain", http.StatusBadRequest)
			return
		}
	}

	// Update the user object
	user := users.Users[userReq.Username]
	if userReq.Displayname != "" {
		user.Displayname = userReq.Displayname
	}
	if userReq.Email != "" {
		user.Email = userReq.Email
	}
	if userReq.Groups != nil {
		user.Groups = userReq.Groups
	}
	users.Users[userReq.Username] = user

	// Convert the users object to YAML string
	usersData, err = yaml.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to marshal users object", http.StatusInternalServerError)
		return
	}

	// Write the updated users.yaml file
	err = os.WriteFile(filePath, usersData, 0644)
	if err != nil {
		http.Error(w, "Failed to write users.yaml file", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "User updated successfully")
}
