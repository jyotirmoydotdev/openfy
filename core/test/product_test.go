package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/jyotirmoydotdev/openfy/db/models"
	web "github.com/jyotirmoydotdev/openfy/internal/web/handlers"
)

type ProductResponses struct {
	Data []models.Product `json:"data"`
}
type ProductResponse struct {
	Data models.Product `json:"data"`
}

var productResponse ProductResponses
var singleProductResponse ProductResponse

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
	req, err := http.NewRequest("POST", server.URL+"/admin/product/new", bytes.NewBuffer(jsonProduct))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
	}
	defer resp.Body.Close()

	if status2 := resp.StatusCode; status2 != http.StatusOK {
		t.Errorf("handler returned wrong staus code: got %v want %v", status2, http.StatusOK)
	}
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
	req, err := http.NewRequest("POST", server.URL+"/admin/product/new", bytes.NewBuffer(jsonProduct))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error making request:%v", err)
	}
	defer resp2.Body.Close()

	if status2 := resp2.StatusCode; status2 != http.StatusOK {
		t.Errorf("handler returned wrong staus code: got %v want %v", status2, http.StatusOK)
	}
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

	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong staus code: got %v want %v", status, http.StatusOK)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Error reading response body:", err)
	}

	if err := json.Unmarshal([]byte(string(body)), &productResponse); err != nil {
		t.Fatalf("Error unmarshalling JSON:%v", err)
		return
	}
}
func TestGetProduct(t *testing.T) {
	product := productResponse.Data[1]
	req, err := http.NewRequest("GET", server.URL+"/admin/product?id="+strconv.FormatUint(uint64(product.ID), 10), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error making request:%v", err)
	}
	defer resp.Body.Close()

	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong staus code: got %v want %v", status, http.StatusOK)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Error reading response body:", err)
	}

	if err := json.Unmarshal([]byte(string(body)), &singleProductResponse); err != nil {
		t.Fatalf("Error unmarshalling JSON:%v", err)
		return
	}
	if !reflect.DeepEqual(singleProductResponse.Data, product) {
		t.Errorf("The Product is not same")
	}
}
func TestUpdateProduct(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	jsonFilePath := filepath.Join(currentDir, "jsonExample", "updateProduct2.json")
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
	id := productResponse.Data[1].ID

	err = json.Unmarshal(Content, &updatedProduct)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	jsonProduct, err := json.Marshal(updatedProduct)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("PUT", server.URL+"/admin/product?id="+strconv.FormatUint(uint64(id), 10), bytes.NewBuffer(jsonProduct))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error making request:%v", err)
	}
	defer resp.Body.Close()

	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong staus code: got %v want %v", status, http.StatusOK)
	}
}

func TestDeleteProductVariant(t *testing.T) {
	product := productResponse.Data[0]
	req, err := http.NewRequest(
		"DELETE",
		server.URL+"/admin/variant?id="+
			strconv.FormatUint(uint64(product.ID), 10)+"&vid="+
			strconv.FormatUint(uint64(product.Variants[0].ID), 10),
		bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error making request:%v", err)
	}
	defer resp.Body.Close()

	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong staus code: got %v want %v", status, http.StatusOK)
	}
}

func TestDeleteProduct(t *testing.T) {
	product := productResponse.Data[0]
	req, err := http.NewRequest("DELETE", server.URL+"/admin/product?id="+strconv.FormatUint(uint64(product.ID), 10), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error making request:%v", err)
	}
	defer resp.Body.Close()

	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong staus code: got %v want %v", status, http.StatusOK)
	}
}

func TestGetAllActiveProduct(t *testing.T) {
	req, err := http.NewRequest("GET", server.URL+"/products", bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error making request:%v", err)
	}
	defer resp.Body.Close()

	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong staus code: got %v want %v", status, http.StatusOK)
	}
}
