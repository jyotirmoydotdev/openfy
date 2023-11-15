package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jyotirmoydotdev/openfy/auth"
)

// Check it a new user can signup or not
// Expected : 200
func TestUserSignup(t *testing.T) {
	newUser := auth.User{
		Username: "testuser",
		Password: "testpassword",
		Email:    "test@example.com",
	}

	jsonUser, err := json.Marshal(newUser)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(server.URL+"/signup", "application/json", bytes.NewBuffer(jsonUser))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong staus code: got %v want %v", status, http.StatusOK)
	}
}

// Check if a new user can login or not
// Expected : 200
func TestUserLogin(t *testing.T) {
	loginCredentials := map[string]string{
		"username": "testuser",
		"password": "testpassword",
	}
	jsonCredentials, err := json.Marshal(loginCredentials)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(server.URL+"/login", "application/json", bytes.NewBuffer(jsonCredentials))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

// Check is same username can signup
// Expected: 400
func TestFailSameUsernaem(t *testing.T) {
	newUser := auth.User{
		Username: "testuser",
		Password: "testpassword2",
		Email:    "test2@example.com",
	}
	jsonUser, err := json.Marshal(newUser)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(server.URL+"/signup", "application/json", bytes.NewBuffer(jsonUser))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if status := resp.StatusCode; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong staus code: got %v want %v", status, http.StatusBadRequest)
	}
}
