package service_test

import (
	"testing"

	"github.com/tchaudhry91/rainbow/service"
)

func TestInMemStore(t *testing.T) {
	store := service.NewInMemStore()

	type TestCase struct {
		Blob string
		Hash string
	}

	cases := []TestCase{
		{"thisisastring", "572642d5581b8b466da59e87bf267ceb7b2afd880b59ed7573edff4d980eb1d5"},
		{"password", "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8"},
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}, // Hash of empty-string https://www.di-mgt.com.au/sha_testvectors.html
	}

	// Add a few values
	for _, c := range cases {
		err := store.Put(c.Blob, c.Hash)
		if err != nil {
			t.Errorf("Error returned while added entry: %s", c.Blob)
		}
	}

	// Retrieve the values
	for _, c := range cases {
		blob, err := store.Get(c.Hash)
		if err != nil {
			t.Errorf("Error while fetching expected entry: %s", c.Hash)
		}
		if blob != c.Blob {
			t.Errorf("Mis-matched Blob for %s, Want=%s, Have=%s", c.Hash, c.Blob, blob)
		}
	}

	// Check for non existent value
	_, err := store.Get("123123")
	if err == nil {
		t.Errorf("Failed to error for a non-existent value")
	}
}
