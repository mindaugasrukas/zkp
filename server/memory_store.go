package server

import (
	"errors"

	"github.com/mindaugasrukas/zkp_example/zkp"
)

var (
	UserExistsError = errors.New("user already exists")
	UserDoesNotExistError = errors.New("user doesn't exist")
)

type (
	Data struct {
		commits zkp.Commits
	}
	InMemoryStore struct {
		store map[zkp.UUID]*Data
	}
)

// NewInMemoryStore returns a new store
func NewInMemoryStore() InMemoryStore {
	return InMemoryStore{
		store: make(map[zkp.UUID]*Data),
	}
}

// Add user to the store
func (m InMemoryStore) Add(user zkp.UUID, commits zkp.Commits) error {
	if _, ok := m.store[user]; ok {
		return UserExistsError
	}
	m.store[user] = &Data{commits: commits}
	return nil
}

// Get user data from the store
func (m InMemoryStore) Get(user zkp.UUID) (*Data, error) {
	data, ok := m.store[user]
	if !ok {
		return nil, UserDoesNotExistError
	}
	return data, nil
}
