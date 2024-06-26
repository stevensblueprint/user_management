package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"

	"user_management/utils"

	"github.com/jordan-wright/email"
	"github.com/knadh/koanf/v2"
	"github.com/redis/go-redis/v9"
)

type RegisterRequest struct {
	Displayname string `json:"displayname"`
	Email       string `json:"email"`
}

type WelcomeTemplateData struct {
	Displayname string
	SignupUrl   string
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request, configFile *koanf.Koanf, redisClient *redis.Client, ctx context.Context) {
	// POST /v1/users/register
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get request body
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

	// Creates token
	token := utils.GenerateSecureToken(10)
	token = utils.HashToken(token)

	// Encrypts token
	secret := configFile.String("SECRET")
	encryptedToken, err := utils.EncryptString([]byte(secret), token)
	if err != nil {
		http.Error(w, "Unable to encrypt token", http.StatusInternalServerError)
		return
	}

	// Stores token in redis
	err = redisClient.SAdd(ctx, "tokenPool", token).Err()
	if err != nil {
		http.Error(w, "Unable to store token", http.StatusInternalServerError)
		return
	}

	encodedDisplayname := strings.ReplaceAll(userDisplayname, " ", "%20")
	baseUrl := configFile.String("BASE_URL")

	url := fmt.Sprintf(baseUrl+"api/v1/users/register?displayname=%s&token=%s", encodedDisplayname, encryptedToken)

	smtpHost := configFile.String("smtp.HOST")
	smtpPort := configFile.String("smtp.PORT")
	smtpUsername := configFile.String("smtp.USERNAME")
	smtpPassword := configFile.String("smtp.PASSWORD")

	htmlString, err := utils.OutputHTMLString("./static/html/welcome.html", WelcomeTemplateData{
		Displayname: userDisplayname,
		SignupUrl:   url,
	})
	if err != nil {
		http.Error(w, "Unable to generate HTML email", http.StatusInternalServerError)
		return
	}

	// Sends registration email to user
	e := &email.Email{
		From:    fmt.Sprintf(`Stevens Blueprint <%s>`, smtpUsername),
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
		HTML: []byte(htmlString),
	}
	_, err = e.AttachFile("./static/logos/logo.png")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to attach logo", http.StatusInternalServerError)
		return
	}

	err = e.Send(smtpHost+":"+smtpPort, smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost))
	if err != nil {
		http.Error(w, "Unable to send email", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "User created successfully")

}
