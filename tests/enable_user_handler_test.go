package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"user_management/handlers"
)

func TestEnableUserHandlerSuccess(t *testing.T) {
	// Reset test_users.yaml file
	err := resetYAMLFile("test_users.yaml")
	if err != nil {
		t.Fatalf("Failed to reset YAML file: %v", err)
	}

	// Construct the URL with query parameter
	req, err := http.NewRequest("POST", "/enableuser?username=existinguser", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.EnableUserRequestHandler(w, r, "test_users.yaml")
	})

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check status code and response
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "User is already enabled\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
