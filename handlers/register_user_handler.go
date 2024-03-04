package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"user_management/utils"

	"github.com/joho/godotenv"
	"github.com/jordan-wright/email"
	"github.com/redis/go-redis/v9"
)

type RegisterRequest struct {
	Displayname string `json:"displayname"`
	Email       string `json:"email"`
}

var REGISTER_PAGE_URL = "http://localhost:8080/api/v1/users/register"

func RegisterUserHandler(w http.ResponseWriter, r *http.Request, client *redis.Client) {
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

	// Check if display name is empty
	userDisplayname := strings.TrimSpace(userReq.Displayname)
	if userDisplayname == "" {
		http.Error(w, "Display name cannot be empty", http.StatusBadRequest)
		return
	}

	// Check if email is empty
	userEmail := strings.TrimSpace(userReq.Email)
	if userEmail == "" {
		http.Error(w, "Email cannot be empty", http.StatusBadRequest)
		return
	}

	err = godotenv.Load(".env")
	if err != nil {
		http.Error(w, "Error loading .env file", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	// Creates token
	token := utils.GenerateSecureToken(10)
	token = utils.HashToken(token)

	// Stores token in redis
	err = client.SAdd(ctx, "tokenPool", token).Err()
	if err != nil {
		http.Error(w, "Unable to store token", http.StatusInternalServerError)
		return
	}

	// Sends registration email to user
	queryDisplayname := strings.ReplaceAll(userDisplayname, " ", "%20")
	url := fmt.Sprintf(REGISTER_PAGE_URL+"?displayname=%s&token=%s", queryDisplayname, token)
	plainAuthUsername := os.Getenv("PLAIN_AUTH_USERNAME")
	plainAuthPassword := os.Getenv("PLAIN_AUTH_PASSWORD")
	e := &email.Email{
		From:    "Stevens Blueprint <test@sitblueprint.com>",
		To:      []string{userEmail},
		Subject: "Blueprint Registration",
		Text: []byte(fmt.Sprintf(`
		Hello %s!
		We are thrilled to welcome you aboard to Blueprint.

		At Blueprint, we're a community of passionate individuals united by a common mission: to amplify impact through technology solutions that drive positive change. As you start your journey with us, you're becoming a vital part of this mission, and we have every confidence that your contributions will be invaluable.

		To complete your registration please complete the following form: %s.

		Once again, welcome to Blueprint! We can't wait to see the positive impact we'll create together.

		Warm regards,
		Blueprint E-Board`, userDisplayname, url)),
	}
	err = e.Send("smtp.gmail.com:587", smtp.PlainAuth("", plainAuthUsername, plainAuthPassword, "smtp.gmail.com")) // Needs email with permission to use gmail
	if err != nil {
		http.Error(w, "Unable to send email", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "User created successfully")

}
