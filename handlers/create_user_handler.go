package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CreateUserRequest struct {
	Displayname string `json:"displayname"`
	Email       string `json:"email"`
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// POST /v1/users/user
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userReq CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if displayname is empty
	if userReq.Displayname == "" {
		http.Error(w, "Displayname cannot be empty", http.StatusBadRequest)
		return
	}

	// Check if email is empty
	if userReq.Email == "" {
		http.Error(w, "Email cannot be empty", http.StatusBadRequest)
		return
	}

	// Check if email belongs to @sitblueprint.com or @stevens.edu
	if !strings.HasSuffix(userReq.Email, "@sitblueprint.com") && !strings.HasSuffix(userReq.Email, "@stevens.edu") {
		http.Error(w, "Email domain must belong to sitblueprint.com or stevens.edu", http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, "User created successfully")

}
