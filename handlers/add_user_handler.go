package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"user_management/models"
	"user_management/utils"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/argon2"
	"gopkg.in/yaml.v2"
)

type AddUserRequest struct {
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Displayname string   `json:"displayname"`
	Email       string   `json:"email"`
	Groups      []string `json:"groups"`
}

/*
AddUserHandler handles the POST /v1/users/user endpoint.
It adds a new user to the users.yaml file. It reads the request body and
token available in the request header. If the token is in the token pool
it will succesfully add the user to the users.yaml file. Else it will
return a forbidden error.
*/
func AddUserHandler(w http.ResponseWriter, r *http.Request, filePath string, secret string, redisClient *redis.Client, ctx context.Context) {
	// POST /v1/users/user
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get request body
	var userReq AddUserRequest
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get authorization header
	encryptedToken := r.Header.Get("Authorization")
	if encryptedToken == "" {
		http.Error(w, "No authorization headers", http.StatusUnauthorized)
		return
	}

	// Parse authorization header
	splitToken := strings.Split(encryptedToken, "Bearer ")
	if len(splitToken) != 2 {
		http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
		return
	}

	// Check if token exists
	encryptedToken = splitToken[1]
	if encryptedToken == "" {
		http.Error(w, "No authorization token found", http.StatusUnauthorized)
		return
	}

	// Decrypt request token
	token, err := utils.DecryptString([]byte(secret), encryptedToken)
	if err != nil {
		http.Error(w, "Unable to decrypt token", http.StatusInternalServerError)
		return
	}

	// Check if decrypted token is in pool of valid tokens
	val, err := redisClient.SIsMember(ctx, "tokenPool", token).Result()
	if err != nil {
		http.Error(w, "Unable to access redis", http.StatusInternalServerError)
		return
	}
	if !val {
		http.Error(w, "Invalid credentials", http.StatusForbidden)
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

	// Check if the username already exists
	if _, exists := users.Users[userReq.Username]; exists {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := hashPassword(userReq.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Create a new user object
	newUser := models.User{
		Disabled:    false,
		Displayname: userReq.Displayname,
		Password:    hashedPassword,
		Email:       userReq.Email,
		Groups:      userReq.Groups,
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

	_, err = redisClient.SRem(ctx, "tokenPool", encryptedToken).Result()
	if err != nil {
		http.Error(w, "Unable to access redis", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "User added successfully")
}

// hashPassword hashes the password using Argon2
func hashPassword(password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Hash the password using Argon2
	hashedPassword := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Encode the hashed password and salt to base64
	encodedHash := base64.RawStdEncoding.EncodeToString(hashedPassword)
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)

	// Concatenate the encoded hash and salt with a separator
	hashed := fmt.Sprintf("$argon2id$v=19$m=65536,t=1,p=4$%s$%s", encodedSalt, encodedHash)

	return hashed, nil
}
