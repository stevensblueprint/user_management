package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"user_management/models"
)

type UserRequest struct {
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Displayname string   `json:"displayname"`
	Email       string   `json:"email"`
	Groups      []string `json:"groups"`
}

func AddUserHandler(w http.ResponseWriter, r *http.Request, filePath string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userReq UserRequest
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	displayname := r.FormValue("displayname")
	email := r.FormValue("email")
	groups := strings.Split(r.FormValue("groups"), ",")

	// Read the existing users.yaml file
	usersFile, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		http.Error(w, "Failed to open users.yaml file", http.StatusInternalServerError)
		return
	}

	var usersData []byte
	usersData, err = io.ReadAll(usersFile)
	if err != nil {
		http.Error(w, "Failed to read users.yaml file", http.StatusInternalServerError)
		return
	}

	defer usersFile.Close()

	var users models.Users
	err = yaml.Unmarshal(usersData, &users)
	if err != nil {
		http.Error(w, "Failed to parse users.yaml file", http.StatusInternalServerError)
		return
	}

	// Check if the username already exists
	if _, exists := users.Users[username]; exists {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	// Create a new user object
	newUser := models.User{
		Disabled:    false,
		Displayname: displayname,
		Password:    password,
		Email:       email,
		Groups:      groups,
	}

	// Add the new user to the users object
	users.Users[username] = newUser

	// Convert the users object to a YAML string
	usersYAML, err := yaml.Marshal(&users)
	if err != nil {
		http.Error(w, "Failed to convert users object to YAML", http.StatusInternalServerError)
		return
	}

	// Reset the file pointer to the beginning of the file
	_, err = usersFile.Seek(0, 0)
	if err != nil {
		http.Error(w, "Failed to seek to the beginning of the file", http.StatusInternalServerError)
		return
	}

	// Write the YAML string to the users.yaml file
	_, err = usersFile.Write(usersYAML)
	if err != nil {
		http.Error(w, "Failed to write users.yaml file", http.StatusInternalServerError)
		return
	}

	// Truncate the file after the new data to remove any old data that might be left after the new data
	err = usersFile.Truncate(int64(len(usersYAML)))
	if err != nil {
		http.Error(w, "Failed to truncate the file", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "User added successfully")
}
