package test

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jyotirmoydotdev/openfy/web"
)

var server *httptest.Server

func startServer() {
	server = httptest.NewServer(web.SetupRouter())
}
func teardown() {
	server.Close()
}

func TestMain(m *testing.M) {
	startServer()

	exitcode := m.Run()
	os.Exit(exitcode)

	defer teardown()
}
