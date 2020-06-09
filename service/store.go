package service

import "fmt"

// Store is an interface defining the operations required from a Rainbow Store
type Store interface {
	Put(blob string, hash string) error
	Get(hash string) (blob string, err error)
}

// InMemStore implements an in-memory map that can function as a Rainbow Store
type InMemStore struct {
	db map[string]string
}

// NewInMemStore instantiates a fresh InMemory Rainbow store
func NewInMemStore() *InMemStore {
	return &InMemStore{db: make(map[string]string)}
}

// Put stores the value of the blob and associates it with the hash as the key
func (store *InMemStore) Put(blob string, hash string) error {
	store.db[hash] = blob
	return nil
}

// Get retrieves the value from the database. Errors if the value is not found.
func (store *InMemStore) Get(hash string) (blob string, err error) {
	if blob, ok := store.db[hash]; ok {
		return blob, nil
	}
	return "", fmt.Errorf("Not Found")
}
