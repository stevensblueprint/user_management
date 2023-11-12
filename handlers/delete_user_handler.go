package handlers

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"user_management/models"
)

func DeleteUserHandler(w http.ResponseWriter, r *http.Request, filePath string) {
	if r.Method != http.MethodDelete { //Method check, if it is not a delete request, return Method not allowed
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//extract the username from the delete request
	username := r.FormValue("username")

	// Read the existing users.yaml file
	usersFile, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		http.Error(w, "Failed to open users.yaml file", http.StatusInternalServerError)
		return
	}
	defer usersFile.Close()
	// open users.yaml file in filePath. If there is an error, response Failed to open users.yaml file.
	//defer is used to close the file after being opened.

	var usersData []byte
	usersData, err = ioutil.ReadAll(usersFile)
	if err != nil {
		http.Error(w, "Failed to read users.yaml file", http.StatusInternalServerError)
		return
	}

	var users models.Users
	err = yaml.Unmarshal(usersData, &users)
	if err != nil {
		http.Error(w, "Failed to parse users.yaml file", http.StatusInternalServerError)
		return
	}

	// Check if the user exists
	if _, exists := users.Users[username]; !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Delete the user
	delete(users.Users, username)

	// Convert the users object to a YAML string
	usersYAML, err := yaml.Marshal(&users)
	if err != nil {
		http.Error(w, "Failed to convert users object to YAML", http.StatusInternalServerError)
		return
	}

	// Reset the file pointer to the beginning of the file
	_, err = usersFile.Seek(0, 0)
	if err != nil {
		http.Error(w, "Failed to seek to the beginning of the file", http.StatusInternalServerError)
		return
	}

	// Write the YAML string to the users.yaml file
	_, err = usersFile.Write(usersYAML)
	if err != nil {
		http.Error(w, "Failed to write users.yaml file", http.StatusInternalServerError)
		return
	}

	// Truncate the file after the new data to remove any old data that might be left after the update
	err = usersFile.Truncate(int64(len(usersYAML)))
	if err != nil {
		http.Error(w, "Failed to truncate the file", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "User deleted successfully")
}
