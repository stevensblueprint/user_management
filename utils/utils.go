package utils

import (
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
    disabled: false
    displayname: Blueprint User 2
    password: existingpassword
    email: user2@blueprint.com
    groups:
      - admin`
	return os.WriteFile(filePath, []byte(initialState), 0644)
}
