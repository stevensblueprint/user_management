package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"user_management/handlers"
	"user_management/utils"

	"github.com/redis/go-redis/v9"
)

func TestAddUserHandlerSuccess(t *testing.T) {
	// Connects to database 1
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})
	ctx := context.Background()

	// Reset the test database
	err := redisClient.Del(ctx, "tokenPool").Err()
	if err != nil {
		t.Fatalf(("Failed to reset database"))
	}

	// Testing variables
	token := "token123"
	secret := "b6fd13cc00dda2a715962dfe6ec32ad0"

	// Store test token
	err = redisClient.SAdd(ctx, "tokenPool", token).Err()
	if err != nil {
		t.Fatalf(("Failed to store test token"))
	}

	// Encrypts test token
	encryptedToken, err := utils.EncryptString([]byte(secret), token)

	// Reset the test_users.yaml file
	err = resetYAMLFile("test_users.yaml")
	if err != nil {
		t.Fatalf("Failed to reset YAML file: %v", err)
	}

	// Setup request body
	userReq := handlers.AddUserRequest{
		Username:    "newuser",
		Password:    "password",
		Displayname: "New User",
		Email:       "example@sitblueprint.com",
		Groups:      []string{"dev", "admin"},
	}
	body, _ := json.Marshal(userReq)

	// Create an HTTP request
	req, err := http.NewRequest("POST", "/v1/users/user", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Sets authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", encryptedToken))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.AddUserHandler(w, r, "test_users.yaml", secret, redisClient, ctx)
	})

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code and response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "User added successfully"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestAddUserHandlerInvalidMethod(t *testing.T) {
	// Connects to database 1
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	ctx := context.Background()

	// Reset the test database
	err := redisClient.Del(ctx, "tokenPool").Err()
	if err != nil {
		t.Fatalf(("Failed to reset database"))
	}

	// Testing variables
	token := "token123"
	secret := "b6fd13cc00dda2a715962dfe6ec32ad0"

	// Store test token
	err = redisClient.SAdd(ctx, "tokenPool", token).Err()
	if err != nil {
		t.Fatalf(("Failed to store test token"))
	}

	// Reset the test_users.yaml file
	err = resetYAMLFile("test_users.yaml")
	if err != nil {
		t.Fatalf("Failed to reset YAML file: %v", err)
	}

	req, err := http.NewRequest("GET", "/v1/users/user", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Sets wrong authorization header
	req.Header.Set("Authorization", "Bearer wrongtoken")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.AddUserHandler(w, r, "test_users.yaml", secret, redisClient, ctx)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}
