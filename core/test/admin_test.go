package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	database "github.com/jyotirmoydotdev/openfy/db/repositories"
	web "github.com/jyotirmoydotdev/openfy/internal/web/handlers"
)

var token string

// Check if Admin can signup or not
// Expected : 200
func TestAdminSignup(t *testing.T) {
	type NewAdminStruct struct {
		Email     string
		Username  string
		FirstName string
		LastName  string
		Password  string
	}
	newAdmin := NewAdminStruct{
		Email:     "test@example.com",
		Username:  "testadmin",
		FirstName: "Test",
		LastName:  "Admin",
		Password:  "testpassword",
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
		Username  string
		Password  string
		Email     string
		FirstName string
		LastName  string
	}
	newAdmin := NewAdminStruct{
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
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	jsonFilePath := filepath.Join(currentDir, "jsonExample", "product1.json")
	file, err := os.Open(jsonFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	Content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var newProduct web.RequestProduct

	err = json.Unmarshal(Content, &newProduct)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
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

func TestUpdateProduct(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	jsonFilePath := filepath.Join(currentDir, "jsonExample", "updateProduct1.json")
	file, err := os.Open(jsonFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	Content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	var updatedProduct web.RequestProduct
	id := database.ProductList[0].ID

	err = json.Unmarshal(Content, &updatedProduct)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	jsonProduct, err := json.Marshal(updatedProduct)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("PUT", server.URL+"/admin/products/"+id, bytes.NewBuffer(jsonProduct))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
	}
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong staus code: got %v want %v", status, http.StatusOK)
	}
	resp.Body.Close()
}

func TestGetAllProducts(t *testing.T) {
	var emptyJson []byte
	req, err := http.NewRequest("GET", server.URL+"/admin/products", bytes.NewBuffer(emptyJson))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
	}
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong staus code: got %v want %v", status, http.StatusOK)
	}
	defer resp.Body.Close()
}
