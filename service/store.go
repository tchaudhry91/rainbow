package service

import (
	"context"
	"fmt"

	redis "github.com/go-redis/redis/v8"
)

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

// RedisStore implements a rainbow store backed by redis
type RedisStore struct {
	client *redis.Client
}

// NewRedisStore instantiates a new redis-backed Rainbow store
func NewRedisStore(addr string, password string, db int) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisStore{
		client: client,
	}
}

// Put sets the blob as the value for the given hash key
func (store *RedisStore) Put(blob string, hash string) error {
	return store.client.Set(context.Background(), hash, blob, 0).Err()
}

// Get retrieves the value from the redis database. Errors if value is not found
func (store *RedisStore) Get(hash string) (blob string, err error) {
	return store.client.Get(context.Background(), hash).Result()
}
