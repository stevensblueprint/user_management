package main

import (
	"os"
)

func resetYAMLFile(filePath string) error {
	initialState := `users:
  existinguser:
    disabled: false
    displayname: Existing User
    password: existingpassword
    email: existinguser@example.com
    groups:
    - group1`
	return os.WriteFile(filePath, []byte(initialState), 0644)
}
