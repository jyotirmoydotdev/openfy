package test

import (
	"database/sql"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

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
	dbInstance, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		return fmt.Errorf("error opening test database: %v", err)
	}
	defer dbInstance.Close()

	statements := []string{
		"DELETE FROM users;",
		"DELETE FROM user_secrets;",
		"DELETE FROM shop_details;",
		"DELETE FROM delivery_addresses;",
		"DELETE FROM sqlite_sequence;",
	}

	for _, statement := range statements {
		_, err := dbInstance.Exec(statement)
		if err != nil {
			return fmt.Errorf("error executing SQL statement: %v", err)
		}
	}

	return nil
}

func TestMain(m *testing.M) {
	err := resetTestDatabase()
	if err != nil {
		fmt.Printf("Error cleaning up test database: %v\n", err)
		os.Exit(1)
	}
	startServer()
	exitcode := m.Run()
	os.Exit(exitcode)
	defer teardown()
}
