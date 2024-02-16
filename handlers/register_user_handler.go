package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"user_management/utils"
)

type RegisterRequest struct {
	Displayname string `json:"displayname"`
	Email       string `json:"email"`
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	// POST /v1/users/register
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userReq RegisterRequest
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
	displayName := map[string]interface{}{"displayName": userReq.Displayname}
	utils.OutputHTML(w, "templates/register_user.html", displayName)

	fmt.Fprint(w, "User created successfully")

}
