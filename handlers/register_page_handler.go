package handlers

import (
	"net/http"
	"strings"
	"user_management/utils"
)

type SignupTemplateData struct {
	Displayname string
}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	// GET /v1/users/register?displayname={displayName}?token={token}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the display name from the query parameter
	displayname := r.URL.Query().Get("displayname")
	if displayname == "" {
		http.Error(w, "Display name is required", http.StatusBadRequest)
		return
	}

	// Extract the token from the query parameter
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	// Generates html
	utils.OutputHTML(w, "./static/html/signup.html", SignupTemplateData{
		Displayname: strings.ReplaceAll(displayname, "%20", " "),
	})
}
