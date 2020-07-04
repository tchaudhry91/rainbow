package service_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tchaudhry91/rainbow/service"
)

func getServer() (server *service.RainbowHTTP) {
	svc := service.NewSHA256RainbowService(service.NewInMemStore())
	return service.NewRainbowHTTP("127.0.0.1:14141", svc)
}

func TestHashHandler(t *testing.T) {
	// This is a method to demonstrate how to test the handlers.
	// Should ideally be a table driven test with multiple requests later.
	server := getServer()
	rr := httptest.NewRecorder()

	// Create a request
	req, err := http.NewRequest("GET", "/hash", nil)
	if err != nil {
		t.Errorf("Failed to create an HTTP Request: %v", err)
		t.FailNow()
	}
	q := req.URL.Query()
	q.Add("str", "thisisastring")
	req.URL.RawQuery = q.Encode()

	server.Handler().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned non OK status: got %d, want %d", rr.Code, http.StatusOK)
		t.FailNow()
	}

	// Check the data
	expectedHash := "572642d5581b8b466da59e87bf267ceb7b2afd880b59ed7573edff4d980eb1d5"
	resp := struct {
		Hash string
		Err  string
	}{}
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Errorf("Error Decoding Response JSON:%v", err)
		t.FailNow()
	}
	if resp.Hash != expectedHash {
		t.Errorf("Wrong hash returned in response. Got - %s Want - %s", resp.Hash, expectedHash)
	}
}
