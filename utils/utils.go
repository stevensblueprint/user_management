package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"html/template"

	"net/http"
	"os"
)

func ResetYAMLFile(filePath string) error {
	initialState := `users:
  user1:
    disabled: false
    displayname: Blueprint User 1
    password: existingpassword
    email: user1@blueprint.com
    groups:
      - admin
      - dev
  user2:
    disabled: true
    displayname: Blueprint User 2
    password: existingpassword
    email: user2@blueprint.com
    groups:
      - admin`
	return os.WriteFile(filePath, []byte(initialState), 0644)
}

func OutputHTML(w http.ResponseWriter, filename string, data interface{}) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, "Failed to parse HTML file", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Failed to execute HTML file", http.StatusInternalServerError)
		return
	}
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func HashToken(input string) string {
	plainText := []byte(input)
	sha256Hash := sha256.Sum256(plainText)
	return hex.EncodeToString(sha256Hash[:])
}
