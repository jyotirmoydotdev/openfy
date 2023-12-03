package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	database "github.com/jyotirmoydotdev/openfy/Database"
)

var UserJWT string

// Check it a new user can signup or not
// Expected : 200
func TestUserSignup(t *testing.T) {
	newUser := map[string]string{
		"email":    "testuser@example.com",
		"password": "testpassword",
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
		"email":    "testuser@example.com",
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
	var reponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&reponse)
	if err != nil {
		t.Errorf("Error decoding JSON response:%v", err)
		return
	}
	UserJWT, _ = reponse["token"].(string)
}

// Check is same username can signup
// Expected: 400
func TestFailSameUsernaem(t *testing.T) {
	newUser := database.User{
		Email:    "testuser@example.com",
		Password: "testpassword2",
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
func TestNthUserSignup(t *testing.T) {
	testNthUser := 1
	for i := 0; i < testNthUser; i++ {
		email := strconv.Itoa(i) + "testuser@example.com"
		newUser := map[string]string{
			"email":    email,
			"password": "testpassword",
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
}

func TestUserPingPong(t *testing.T) {
	// Create a request with the correct endpoint
	req, err := http.NewRequest("GET", server.URL+"/api/ping", nil)
	if err != nil {
		t.Fatal("Error creating request:", err)
	}

	// Set the Authorization header with the UserJWT
	req.Header.Set("Authorization", "Bearer "+UserJWT)

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("Error making request:", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Decode the JSON response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Errorf("Error decoding JSON response: %v", err)
		return
	}

	// Check if the "message" field exists in the response
	message, ok := response["message"].(string)
	if !ok {
		t.Error("Expected 'message' field in response, but it was not found")
	} else {
		fmt.Println(message)
	}
}
