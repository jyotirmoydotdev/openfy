package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

var token string

// Check if Admin can signup or not
// Expected : 200
func TestAdminSignup(t *testing.T) {
	type NewAdminStruct struct {
		Email        string
		Customername string
		FirstName    string
		LastName     string
		Password     string
	}
	newAdmin := NewAdminStruct{
		Email:        "test@example.com",
		Customername: "testadmin",
		FirstName:    "Test",
		LastName:     "Admin",
		Password:     "testpassword",
	}
	jsonAdmin, err := json.Marshal(newAdmin)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(server.URL+"/admin/signup", "application/json", bytes.NewBuffer(jsonAdmin))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong staus code: got %v want %v", status, http.StatusOK)
	}
}

// Check Admin can login or not
// Expected : 200
func TestAdminLogin(t *testing.T) {
	loginCredentials := map[string]string{
		"customername": "testadmin",
		"password":     "testpassword",
	}
	jsonCredentials, err := json.Marshal(loginCredentials)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(server.URL+"/admin/login", "application/json", bytes.NewBuffer(jsonCredentials))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong staus code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Errorf("Error decoding JSON response:%v", err)
		return
	}
	var ok bool
	token, ok = response["token"].(string)
	if !ok {
		fmt.Println("Token not found in the reponse")
		return
	}
}

// Check if second new admin can signup or not
// Expected : 403
func TestFailAdminSignup(t *testing.T) {
	type NewAdminStruct struct {
		Customername string
		Password     string
		Email        string
		FirstName    string
		LastName     string
	}
	newAdmin := NewAdminStruct{
		Customername: "testadmin1",
		Password:     "testpassword",
		Email:        "test@example.com",
		FirstName:    "Test",
		LastName:     "Admin",
	}
	jsonAdmin, err := json.Marshal(newAdmin)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(server.URL+"/admin/signup", "application/json", bytes.NewBuffer(jsonAdmin))
	if err != nil {
		t.Fatal(err)
	}
	if status := resp.StatusCode; status != http.StatusForbidden {
		t.Errorf("handler returned wrong staus code: got %v want %v", status, http.StatusForbidden)
	}
	defer resp.Body.Close()
}
