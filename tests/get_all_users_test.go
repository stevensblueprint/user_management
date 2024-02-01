package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"user_management/handlers"
)

func TestGetAllUsersHandlerSuccess(t *testing.T) {
	// Reset the test_users.yaml file
	err := resetYAMLFile("test_users.yaml")
	if err != nil {
		t.Fatalf("Failed to reset YAML file: %v", err)
	}

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAllUsersHandler(w, r, "test_users.yaml")
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if the response body is the expected JSON (you'll need to define what you expect)
	expected := `{"Users":{"existinguser":{"Disabled":false,"Displayname":"Existing User","Password":"existingpassword","Email":"existinguser@example.com","Groups":["group1"]}}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetAllUsersHandlerInvalidMethod(t *testing.T) {
	req, err := http.NewRequest("POST", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAllUsersHandler(w, r, "test_users.yaml")
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestGetAllUsersHandlerFileReadError(t *testing.T) {
	// Provide an invalid file path to simulate file read error
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAllUsersHandler(w, r, "invalid_path.yaml")
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}
