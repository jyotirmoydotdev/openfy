package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/jyotirmoydotdev/openfy/db/models"
)

var CustomerJWT string

// Check it a new customer can signup or not
// Expected : 200
func TestCustomerSignup(t *testing.T) {
	newCustomer := map[string]string{
		"email":     "testcustomer@example.com",
		"password":  "testpassword",
		"firstname": "Jyotirmoy",
		"lastname":  "Barman",
	}

	jsonCustomer, err := json.Marshal(newCustomer)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(server.URL+"/signup", "application/json", bytes.NewBuffer(jsonCustomer))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong staus code: got %v want %v", status, http.StatusOK)
	}
}

// Check if a new customer can login or not
// Expected : 200
func TestCustomerLogin(t *testing.T) {
	loginCredentials := map[string]string{
		"email":    "testcustomer@example.com",
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
	CustomerJWTfetch, ok := reponse["token"].(string)
	if !ok {
		t.Errorf("Something went wrrong while fetching token from the reponse")
	}
	CustomerJWT = CustomerJWTfetch
}
func TestCustomerLogin2(t *testing.T) {
	loginCredentials := map[string]string{
		"email":    "testcustomer@example.com",
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
	CustomerJWTfetch, ok := reponse["token"].(string)
	if !ok {
		t.Errorf("Something went wrrong while fetching token from the reponse")
	}
	CustomerJWT = CustomerJWTfetch
}

// Check is same username can signup
// Expected: 400
func TestFailSameUsername(t *testing.T) {
	newCustomer := models.Customer{
		Email:    "testcustomer@example.com",
		Password: "testpassword2",
	}
	jsonCustomer, err := json.Marshal(newCustomer)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(server.URL+"/signup", "application/json", bytes.NewBuffer(jsonCustomer))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if status := resp.StatusCode; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong staus code: got %v want %v", status, http.StatusBadRequest)
	}
}
func TestNthCustomerSignup(t *testing.T) {
	testNthCustomer := 10
	for i := 0; i < testNthCustomer; i++ {
		email := strconv.Itoa(i) + "testcustomer@example.com"
		newCustomer := map[string]string{
			"email":     email,
			"password":  "testpassword",
			"firstname": strconv.Itoa(i) + "Jyotirmoy",
			"lastname":  strconv.Itoa(i) + "Barman",
		}
		jsonCustomer, err := json.Marshal(newCustomer)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.Post(server.URL+"/signup", "application/json", bytes.NewBuffer(jsonCustomer))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if status := resp.StatusCode; status != http.StatusOK {
			t.Errorf("handler returned wrong staus code: got %v want %v", status, http.StatusOK)
		}
	}
}
func TestNthCustomerLogin(t *testing.T) {
	testNthCustomer := 10
	for i := 0; i < testNthCustomer; i++ {
		email := strconv.Itoa(i) + "testcustomer@example.com"
		newCustomer := map[string]string{
			"email":    email,
			"password": "testpassword",
		}
		jsonCustomer, err := json.Marshal(newCustomer)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.Post(server.URL+"/login", "application/json", bytes.NewBuffer(jsonCustomer))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if status := resp.StatusCode; status != http.StatusOK {
			t.Errorf("handler returned wrong staus code: got %v want %v", status, http.StatusOK)
		}
	}
}

func TestCustomerPingPong(t *testing.T) {
	// Create a request with the correct endpoint
	req, err := http.NewRequest("GET", server.URL+"/customer/ping", nil)
	if err != nil {
		t.Fatal("Error creating request:", err)
	}

	// Set the Authorization header with the CustomerJWT
	req.Header.Set("Authorization", "Bearer "+CustomerJWT)

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
	}
	if message != "pong" {
		t.Error("Message does not match")
	}
}
