package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"user_management/handlers"
)

func TestUpdateUserHandlerSuccess(t *testing.T) {
	// Reset the test_users.yaml file
	err := resetYAMLFile("test_users.yaml")
	if err != nil {
		t.Fatalf("Failed to reset YAML file: %v", err)
	}

	// Setup request body
	userReq := handlers.UpdateUserRequest{
		Username:    "existinguser",
		Displayname: "Updated User",
	}
	body, _ := json.Marshal(userReq)

	// Create an HTTP request
	req, err := http.NewRequest("PUT", "/updateuser", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateUserHandler(w, r, "test_users.yaml")
	})

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code and response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "User updated successfully"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
