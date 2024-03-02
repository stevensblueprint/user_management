package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"user_management/models"

	"gopkg.in/yaml.v2"
)

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request, filePath string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

	// Convert the users object to JSON
	usersJSON, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to convert users to JSON", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(usersJSON)
}
