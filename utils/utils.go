package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"html/template"
	"io"
	"strings"

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

func OutputHTMLString(filename string, data interface{}) (string, error) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		return "", err
	}

	builder := &strings.Builder{}

	if err := t.Execute(builder, data); err != nil {
		return "", err
	}

	return builder.String(), nil
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

func EncryptString(key []byte, text string) (string, error) {
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptString(key []byte, encryptedString string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(encryptedString)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
