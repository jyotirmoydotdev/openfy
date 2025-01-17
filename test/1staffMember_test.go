package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

var token string

// Check if StaffMember can signup or not
// Expected : 200
func TestStaffMemberSignup(t *testing.T) {
	type NewStaffMemberStruct struct {
		Email     string
		Username  string
		FirstName string
		LastName  string
		Password  string
	}
	newStaffMember := NewStaffMemberStruct{
		Email:     "test@example.com",
		Username:  "teststaffMember",
		FirstName: "Test",
		LastName:  "StaffMember",
		Password:  "testpassword",
	}
	jsonStaffMember, err := json.Marshal(newStaffMember)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(server.URL+"/staffMember/signup", "application/json", bytes.NewBuffer(jsonStaffMember))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong staus code: got %v want %v", status, http.StatusOK)
	}
}

// Check StaffMember can login or not
// Expected : 200
func TestStaffMemberLogin(t *testing.T) {
	loginCredentials := map[string]string{
		"username": "teststaffMember",
		"password": "testpassword",
	}
	jsonCredentials, err := json.Marshal(loginCredentials)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(server.URL+"/staffMember/login", "application/json", bytes.NewBuffer(jsonCredentials))
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

// Check if second new staffMember can signup or not
// Expected : 403
func TestFailStaffMemberSignup(t *testing.T) {
	type NewStaffMemberStruct struct {
		Username  string
		Password  string
		Email     string
		FirstName string
		LastName  string
	}
	newStaffMember := NewStaffMemberStruct{
		Username:  "teststaffMember1",
		Password:  "testpassword",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "StaffMember",
	}
	jsonStaffMember, err := json.Marshal(newStaffMember)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(server.URL+"/staffMember/signup", "application/json", bytes.NewBuffer(jsonStaffMember))
	if err != nil {
		t.Fatal(err)
	}
	if status := resp.StatusCode; status != http.StatusForbidden {
		t.Errorf("handler returned wrong staus code: got %v want %v", status, http.StatusForbidden)
	}
	defer resp.Body.Close()
}
