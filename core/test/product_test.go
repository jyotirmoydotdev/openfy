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

	web "github.com/jyotirmoydotdev/openfy/internal/web/handlers"
)

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
func TestAddProduct2(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	jsonFilePath := filepath.Join(currentDir, "jsonExample", "product2.json")
	file, err := os.Open(jsonFilePath)
	if err != nil {
		t.Fatalf("Error opening file:%v", err)
		return
	}
	defer file.Close()

	Content, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("Error reading file:%v", err)
		return
	}

	var newProduct web.RequestProduct

	err = json.Unmarshal(Content, &newProduct)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON:%v", err)
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
		t.Fatalf("Error making request:%v", err)
	}
	if status2 := resp2.StatusCode; status2 != http.StatusOK {
		t.Errorf("handler returned wrong staus code: got %v want %v", status2, http.StatusOK)
	}
	resp2.Body.Close()
}
func TestGetAllProduct(t *testing.T) {
	req, err := http.NewRequest("GET", server.URL+"/admin/products", bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error making request:%v", err)
	}
	defer resp.Body.Close()

	if status2 := resp.StatusCode; status2 != http.StatusOK {
		t.Errorf("handler returned wrong staus code: got %v want %v", status2, http.StatusOK)
	}
}
