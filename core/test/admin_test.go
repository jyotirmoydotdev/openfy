package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jyotirmoydotdev/openfy/auth"
	"github.com/jyotirmoydotdev/openfy/web"
)

var token string

// Check if Admin can signup or not
// Expected : 200
func TestAdminSignup(t *testing.T) {
	newAdmin := auth.Admin{
		Username:  "testadmin",
		Password:  "testpassword",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "Admin",
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
		"username": "testadmin",
		"password": "testpassword",
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
		fmt.Println("Error decoding JSON response:", err)
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
	newAdmin := auth.Admin{
		Username:  "testadmin1",
		Password:  "testpassword",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "Admin",
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

func TestAddProduct(t *testing.T) {
	newProduct := web.Product{
		Title:            "Sample Product",
		Description:      "This is a sample product description.",
		Media:            []string{"https://example.com/image1.jpg", "https://example.com/image2.jpg"},
		Price:            100,
		Compare_At_Price: 120,
		Tax:              true,
		Cost_Per_Item:    80,
	}

	jsonProduct, err := json.Marshal(newProduct)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", server.URL+"/admin/products/new", bytes.NewBuffer(jsonProduct))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
	}
	if status2 := resp2.StatusCode; status2 != http.StatusOK {
		t.Errorf("handler returned wrong staus code: got %v want %v", status2, http.StatusOK)
	}
	resp2.Body.Close()
}
