package service

import (
	"crypto/sha256"
	"encoding/hex"
)

// Hasher is any type that can return a hash of string
type Hasher interface {
	Hash(blob string) (hashed string, err error)
}

// HashReverser is any type that can reverse a hash and return the original blob
type HashReverser interface {
	HashReverse(hashed string) (blob string, err error)
}

// RainbowService is a service to compute hashes and lookup reverse hashes
type RainbowService interface {
	Hasher
	HashReverser
}

// SHA256RainbowService is a SHA256 implementation for the RainbowService
type SHA256RainbowService struct {
	store Store
}

// NewSHA256RainbowService instantiates a SHA256RainbowService
func NewSHA256RainbowService(store Store) *SHA256RainbowService {
	return &SHA256RainbowService{store}
}

// Hash returns the SHA256 sum for the given string
func (svc *SHA256RainbowService) Hash(blob string) (hashed string, err error) {
	sumA := sha256.Sum256([]byte(blob))
	hashed = hex.EncodeToString(sumA[:])
	err = svc.store.Put(blob, hashed)
	return
}

// HashReverse looks up the original string for the given hash
func (svc *SHA256RainbowService) HashReverse(hashed string) (blob string, err error) {
	return svc.store.Get(hashed)
}
