package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"

	"user_management/handlers"
)

type UserRequest struct {
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Displayname string   `json:"displayname"`
	Email       string   `json:"email"`
	Groups      []string `json:"groups"`
}

func readFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func TestAddUserHandler(t *testing.T) {
	// Create a temporary users.yaml file for testing
	tempFile, err := os.CreateTemp("", "users.yaml")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer tempFile.Close()

	// Write initial data to the temporary file
	initialData := `
users:
  existinguser:
    disabled: false
    displayname: Existing User
    password: existingpassword
    email: existinguser@example.com
    groups:
      - group1
`
	err = os.WriteFile(tempFile.Name(), []byte(initialData), 0644)
	if err != nil {
		t.Fatalf("Failed to write initial data to temporary file: %v", err)
	}

	// Create a test request JSON
	userReq := UserRequest{
		Username:    "newuser",
		Password:    "newpassword",
		Displayname: "New User",
		Email:       "newuser@example.com",
		Groups:      []string{"group1", "group2"},
	}

	jsonData, err := json.Marshal(userReq)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "/adduser", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the response
	res := httptest.NewRecorder()

	// Call the addUserHandler function with the test request and response recorder
	handlers.AddUserHandler(res, req, tempFile.Name())

	// Check the response status code
	if res.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, res.Code)
	}

	// Read the updated users.yaml file
	usersData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read users.yaml file: %v", err)
	}

	// Assert that the new user is added to the users.yaml file
	expectedData, err := readFile("test_users.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var usersDataObj, expectedDataObj interface{}
	if err := yaml.Unmarshal(usersData, &usersDataObj); err != nil {
		t.Fatal(err)
	}
	if err := yaml.Unmarshal(expectedData, &expectedDataObj); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(usersDataObj, expectedDataObj) {
		t.Errorf("Expected users.yaml data:\n%v\n\nBut got:\n%v", expectedDataObj, usersDataObj)
	}
}
