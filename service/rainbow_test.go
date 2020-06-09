package service_test

import (
	"testing"

	"github.com/tchaudhry91/rainbow/service"
)

// getServices creates a simple service to test against
func getService() (svc service.RainbowService, err error) {
	svc = service.NewSHA256RainbowService(service.NewInMemStore())
	return
}

func TestServiceHash(t *testing.T) {
	svc, err := getService()
	if err != nil {
		t.Errorf("Failed to get sample test service")
		t.FailNow()
	}

	type TestCase struct {
		Blob       string
		WantedHash string
	}

	cases := []TestCase{
		{"thisisastring", "572642d5581b8b466da59e87bf267ceb7b2afd880b59ed7573edff4d980eb1d5"},
		{"password", "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8"},
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}, // Hash of empty-string https://www.di-mgt.com.au/sha_testvectors.html
	}

	for _, testCase := range cases {
		receivedHash, err := svc.Hash(testCase.Blob)
		if err != nil {
			t.Errorf("Failed to insert the blob %s because %v:", testCase.Blob, err)
		}
		if testCase.WantedHash != receivedHash {
			t.Errorf("Failed for blob - %s, Wanted: %s, Received: %s", testCase.Blob, testCase.WantedHash, receivedHash)
		}
	}
}

func TestServiceHashReverse(t *testing.T) {
	svc, err := getService()
	if err != nil {
		t.Errorf("Failed to get sample test service")
		t.FailNow()
	}

	type TestCase struct {
		Hash       string
		WantedBlob string
	}

	cases := []TestCase{
		{"572642d5581b8b466da59e87bf267ceb7b2afd880b59ed7573edff4d980eb1d5", "thisisastring"},
		{"5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8", "password"},
		{"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", ""}, // Hash of empty-string https://www.di-mgt.com.au/sha_testvectors.html
	}

	for _, testCase := range cases {
		// Put the key in the database first to ensure availability
		_, err := svc.Hash(testCase.WantedBlob)
		if err != nil {
			t.Errorf("Failed to Hash :%s", err)
		}
		receivedBlob, err := svc.HashReverse(testCase.Hash)
		if err != nil {
			t.Errorf("Failed to calculate reverse hash for %s: %v", testCase.Hash, err)
		}
		if testCase.WantedBlob != receivedBlob {
			t.Errorf("Failed for Hash - %s, Wanted: %s, Received: %s", testCase.Hash, testCase.WantedBlob, receivedBlob)
		}
	}
}
