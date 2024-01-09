package test

import (
	"database/sql"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	db "github.com/jyotirmoydotdev/openfy/database"
	"github.com/jyotirmoydotdev/openfy/internal/web"
)

var server *httptest.Server

func startServer() {
	server = httptest.NewServer(web.SetupRouter())
}
func teardown() {
	server.Close()
}
func resetTestDatabase() error {
	dbInstance, err := sql.Open("sqlite3", "./database/databaseCustomer.db")
	if err != nil {
		return fmt.Errorf("error opening test database: %v", err)
	}
	defer dbInstance.Close()

	statements := []string{
		"DELETE FROM customers;",
		"DELETE FROM customer_secrets;",
		"DELETE FROM shop_details;",
		"DELETE FROM delivery_addresses;",
		"DELETE FROM sqlite_sequence;",
		"DELETE FROM customer_tokens;",
	}

	for _, statement := range statements {
		_, err := dbInstance.Exec(statement)
		if err != nil {
			return fmt.Errorf("error executing SQL statement databaseCustomer: %v", err)
		}
	}

	staffMemberDBInstance, err := sql.Open("sqlite3", "./database/databaseStaffMember.db")
	if err != nil {
		return fmt.Errorf("error opening test databaseStaffMember: %v", err)
	}
	defer staffMemberDBInstance.Close()

	statements = []string{
		"DELETE FROM staff_member_secrets;",
		"DELETE FROM staff_members;",
	}

	for _, statement := range statements {
		_, err := staffMemberDBInstance.Exec(statement)
		if err != nil {
			return fmt.Errorf("error executing SQL statement of databaseStaffMember: %v", err)
		}
	}

	productdbInstance, err := sql.Open("sqlite3", "./database/databaseProduct.db")
	if err != nil {
		return fmt.Errorf("error opening test database databaseProduct: %v", err)
	}
	defer productdbInstance.Close()
	statements = []string{
		"DELETE FROM products;",
		"DELETE FROM options;",
		"DELETE FROM variants;",
		"DELETE FROM selected_options;",
	}
	for _, statement := range statements {
		_, err := productdbInstance.Exec(statement)
		if err != nil {
			return fmt.Errorf("error executing SQL statement: %v", err)
		}
	}
	return nil
}
func TestMain(m *testing.M) {
	err := db.InitializeDatabases()
	if err != nil {
		fmt.Printf("Error cleaning up test database: %v\n", err)
		os.Exit(1)
	}
	err = resetTestDatabase()
	if err != nil {
		fmt.Printf("Error cleaning up test database: %v\n", err)
		os.Exit(1)
	}
	startServer()
	exitcode := m.Run()
	os.Exit(exitcode)
	defer teardown()
}
