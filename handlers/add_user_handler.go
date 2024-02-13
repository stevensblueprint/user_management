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

type AddUserRequest struct {
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Displayname string   `json:"displayname"`
	Email       string   `json:"email"`
	Groups      []string `json:"groups"`
}

func AddUserHandler(w http.ResponseWriter, r *http.Request, filePath string) {
	// POST /v1/users/user
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userReq AddUserRequest
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

	// Check if the username already exists
	if _, exists := users.Users[userReq.Username]; exists {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	// Create a new user object
	newUser := models.User{
		Disabled:    false,
		Displayname: userReq.Displayname,
		Password:    userReq.Password,
		Email:       userReq.Email,
		Groups:      userReq.Groups,
	}

	// Check if groups belongs to [admin, dev]
	for _, group := range userReq.Groups {
		if group != "admin" && group != "dev" {
			http.Error(w, "Invalid group", http.StatusBadRequest)
			return
		}
	}

	// Check if email belongs to @sitblueprint.com or @stevens.edu
	if !strings.HasSuffix(userReq.Email, "@sitblueprint.com") && !strings.HasSuffix(userReq.Email, "@stevens.edu") {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	// Add the new user to the users object
	users.Users[userReq.Username] = newUser

	// Convert the users object to a YAML string
	usersYAML, err := yaml.Marshal(&users)
	if err != nil {
		http.Error(w, "Failed to convert users object to YAML", http.StatusInternalServerError)
		return
	}

	// Write the YAML string to the users.yaml file
	err = os.WriteFile(filePath, usersYAML, 0644)
	if err != nil {
		http.Error(w, "Failed to write users.yaml file", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "User added successfully")
}
